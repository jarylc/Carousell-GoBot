package main

import (
	"carousell-gobot/carousell"
	"carousell-gobot/data/config"
	"carousell-gobot/data/state"
	"carousell-gobot/messaging"
	"flag"
	"log"
	"os"
)

var debug = os.Getenv("DEBUG") == "1"

//nolint:funlen,gocognit
func main() {
	var configPath string
	var statePath string
	flag.StringVar(&configPath, "c", "config.yaml", "path to config file")
	flag.StringVar(&statePath, "s", "state.json", "path to state file")
	flag.Parse()

	if debug {
		log.Println("DEBUG MODE ACTIVE!")
	}

	log.Println("loading config...")
	err := config.Load(configPath)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("config loaded")

	log.Println("loading state...")
	err = state.Load(statePath)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("state loaded")

	log.Println("loading forwarders...")
	messaging.LoadForwarders()
	log.Println("forwarders loaded")

	log.Println("initiating reminders system")
	carousell.InitReminders()
	log.Println("reminder system initiated")

	log.Println("connecting to chat...")
	_ = carousell.Connect()
}
