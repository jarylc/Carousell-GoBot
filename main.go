package main

import (
	"carousell-gobot/carousell"
	"carousell-gobot/data/config"
	"carousell-gobot/data/state"
	"carousell-gobot/forwarders"
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

	log.Println("Loading config...")
	err := config.Load(configPath)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Config loaded")

	log.Println("Loading state...")
	err = state.Load(statePath)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("State loaded")

	log.Println("Loading forwarders...")
	forwarders.LoadForwarders()
	log.Println("Forwarders loaded")

	log.Println("Initiating reminders system")
	carousell.InitReminders()
	log.Println("Reminder system initiated")

	log.Println("Connecting to chat...")
	_ = carousell.Connect()
}
