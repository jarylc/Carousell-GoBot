package carousell

import (
	"carousell-gobot/data/config"
	"testing"
)

// get Carousell user ID
func TestGetUserID(t *testing.T) {
	config.Config.Carousell.Cookie = "_t=t%3D1671458909706%26u%3D344194;"

	userID, err := getUserIDFromCacheOrCookie()
	if err != nil {
		t.Error(err)
	}
	t.Log(userID + " vs 344194")
	if userID != "344194" {
		t.Error("Invalid UserID")
	}
}