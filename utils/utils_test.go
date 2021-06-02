package utils

import (
	"testing"
)

var positive = map[string]float64{
	"0.01":             0.01,
	"able to do $2.5?": 2.5,
	"can 10?":          10,
	"10 can?":          10,
	"offer 10":         10,
	"5?":               5.0,
	"please do 5.50":   5.50,
	"5.50 pls":         5.50,
	"fast deal 99.99":  99.99,
	"10 deal today":    10,
	"20 quick deal":    20,
	"fast deal 420.69": 420.69,
	"sell me at 50":    50,
	"please 0.1":       0.1,
}

var negative = []string{
	"can we deal at 5pm",
	"deal 10am",
	"can we meet at block 253",
	"612345",
	"deal at 654321",
	"can we deal 5pm tomorrow",
	"6am deal",
	"deal 11 tomorrow night",
	"deal 6 at night",
	"deal 7 in the evening",
	"deal at 612435 faster",
	"please meet me at 10 tonight",
}

// get price from message
func TestGetPriceFromMessage(t *testing.T) {
	// test positive
	for msg, actual := range positive {
		price, err := GetPriceFromMessage(msg)
		if err != nil {
			t.Error(err)
		}
		if price != actual {
			t.Errorf("Mismatch `%s` : %.02f", msg, price)
		}
	}

	// test negative
	for _, msg := range negative {
		price, err := GetPriceFromMessage(msg)
		if err != nil {
			t.Error(err)
		}
		if price > 0 {
			t.Errorf("Negative case hit `%s`", msg)
		}
	}
}
