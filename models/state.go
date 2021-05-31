package models

import (
	"time"
)

type State struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	DealOn       time.Time `json:"deal_on"`
	LastReceived string    `json:"last_received"`
	LastSent     string    `json:"last_sent"`
	LastActivity time.Time `json:"last_activity"`
}

func NewState(id string) *State {
	return &State{
		ID:           id,
		Name:         "",
		Price:        0,
		DealOn:       time.Time{},
		LastReceived: "",
		LastSent:     "",
		LastActivity: time.Time{},
	}
}
