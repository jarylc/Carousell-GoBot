package chrono

import (
	"testing"
	"time"
)

var tests = map[string]time.Duration{
	"1 hour later": time.Hour * 1,
	"tomorrow":     time.Hour * 24,
}

func TestParseDate(t *testing.T) {
	chrono, err := New()
	if err != nil {
		t.Error(err)
	}

	now := time.Now()

	for text, duration := range tests {
		date, err := chrono.ParseDate(text, now)
		if err != nil {
			t.Error(err)
		}

		expected := now.Add(duration).Truncate(1 * time.Second)
		if !date.Equal(expected) {
			t.Errorf("%s != %s", date.String(), expected.String())
		}
	}
}
