package models

import (
	"time"
)

type State struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	DealOn       time.Time `json:"deal_on"`
	LastResponse string    `json:"last_response"`
	LastReply    string    `json:"last_reply"`
	LastActivity time.Time `json:"last_activity"`
}

func NewState() *State {
	return &State{
		ID:           "",
		Name:         "",
		Price:        0,
		DealOn:       time.Time{},
		LastResponse: "",
		LastReply:    "",
		LastActivity: time.Time{},
	}
}
