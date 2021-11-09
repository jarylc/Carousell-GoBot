package carousell

import (
	"carousell-gobot/data/config"
	"testing"
)

// get Carousell user ID
func TestGetUserID(t *testing.T) {
	config.Config.Carousell.Cookie = "jwt=stubbed.eyJpZCI6IjM0NDE5NCIsImlzcyI6ImMiLCJpc3N1ZWRfYXQiOjAsInNlY3JldCI6IiIsInVzZXIiOiJqYXJ5bGMifQo.stubbed;"
	userID, err := getUserIDFromCacheOrCookie()
	if err != nil {
		t.Error(err)
	}
	if userID != "344194" {
		t.Error(userID + " vs 344194")
	}
}
