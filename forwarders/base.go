package forwarders

import (
	"carousell-gobot/data/config"
	"log"
)

// Forwarders - contains all the forwarders of the program
var Forwarders []Forwarder

// Forwarder - Base forwarder interface
type Forwarder interface {
	SendMessage(text string) error
	Escape(str string) string
}

// LoadForwarders - load all forwarders from configuration
func LoadForwarders() {
	if config.Config.Forwarders == nil {
		return
	}
	for _, forwarder := range config.Config.Forwarders {
		switch forwarder.Type {
		case "telegram":
			telegram := Telegram{
				Token:  forwarder.Token,
				ChatID: forwarder.ChatID,
			}
			Forwarders = append(Forwarders, telegram)
		}
		log.Printf("\t1x %s loaded\n", forwarder.Type)
	}
}
