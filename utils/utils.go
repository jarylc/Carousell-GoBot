package utils

import (
	"carousell-gobot/constants"
	"carousell-gobot/data/config"
	"context"
	"encoding/json"
	"errors"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/dlclark/regexp2"
	"github.com/jarylc/go-chromedpproxy"
	"regexp"
	"strconv"
	"time"
)

func HTTPGet(url string, out interface{}) error {
	targetID, err := chromedpproxy.NewTab("about://newtab")
	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	defer chromedpproxy.CloseTarget(targetID)
	ctx := chromedpproxy.GetTarget(targetID)

	type result struct {
		Response string
		Error    error
	}
	ch := make(chan result, 1)
	go func() {
		err = chromedp.Run(ctx, chromedp.Tasks{
			network.SetExtraHTTPHeaders(map[string]interface{}{
				"Cookie": config.Config.Carousell.Cookie,
			}),
			chromedp.Navigate(url),
			chromedp.WaitNotPresent(`//div[contains(text(), 'Ray ID')]`),
			chromedp.WaitVisible(`pre`),
			chromedp.ActionFunc(func(ctx context.Context) error {
				node, err := dom.GetDocument().Do(ctx)
				if err != nil {
					return err
				}
				response, err := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
				if err != nil {
					return err
				}
				ch <- result{Response: response, Error: nil}
				return nil
			}),
		})
		if err != nil {
			ch <- result{Response: "", Error: err}
		}
	}()

	res := <-ch
	if res.Error != nil {
		return res.Error
	}

	re := regexp.MustCompile(`\{.+}`)
	resp := re.FindStringSubmatch(res.Response)[0]
	err = json.Unmarshal([]byte(resp), out)
	if err != nil {
		return err
	}

	return nil
}

// GetPriceFromMessage - get price from message, -1 means not detected
func GetPriceFromMessage(msg string) (float64, error) {
	r, err := regexp2.Compile(constants.PRICE_EXPRESSION, regexp2.IgnoreCase)
	if err != nil {
		return 0, err
	}

	priceStr, err := r.FindStringMatch(msg)
	if err != nil {
		return 0, err
	}

	price := 0.0
	if priceStr != nil {
		price, err = strconv.ParseFloat(priceStr.String(), 64)
		if err != nil {
			return 0, err
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
