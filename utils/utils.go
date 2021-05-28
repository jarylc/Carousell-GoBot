package utils

import (
	"carousell-gobot/constants"
	"carousell-gobot/data/config"
	"encoding/json"
	"github.com/dlclark/regexp2"
	"io"
	"net/http"
	"strconv"
	"time"
)

// HTTP functions
var client = &http.Client{}

func HTTPGet(url string, out interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Cookie", config.Config.Carousell.Cookie)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(raw, out)
	if err != nil {
		return err
	}

	return nil
}

// GetPriceFromMessage - get price from message, -1 means not detected
func GetPriceFromMessage(msg string) (float64, error) {
	r, err := regexp2.Compile(constants.PRICE_EXPRESSION, regexp2.IgnoreCase)
	if err != nil {
		return -1, err
	}

	priceStr, err := r.FindStringMatch(msg)
	if err != nil {
		return -1, err
	}

	price := -1.0
	if priceStr != nil {
		price, err = strconv.ParseFloat(priceStr.String(), 64)
		if err != nil {
			return -1, err
		}
	}

	return price, nil
}

// GetEpochString - get epoch as string
func GetEpochString() string {
	return strconv.FormatInt(GetEpoch(), 10)
}

// GetEpoch - get epoch
func GetEpoch() int64 {
	return time.Now().Unix()
}
