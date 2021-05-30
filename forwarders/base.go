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
		var instance Forwarder = nil
		switch forwarder.Type {
		case "telegram":
			instance = Telegram{
				Token:  forwarder.Token,
				ChatID: forwarder.ChatID,
			}
		case "discord":
			instance = Discord{
				WebhookURL: forwarder.WebhookURL,
			}
		case "slack":
			instance = Slack{
				WebhookURL: forwarder.WebhookURL,
			}
		default:
			log.Printf("Skipping invalid forwarder type `%s`\n", forwarder.Type)
		}
		if instance != nil {
			Forwarders = append(Forwarders, instance)
			log.Printf("\t1x %s loaded\n", forwarder.Type)
		}
	}
}
