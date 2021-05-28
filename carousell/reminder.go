package carousell

import (
	"carousell-gobot/data/config"
	"carousell-gobot/data/state"
	"carousell-gobot/forwarders"
	"carousell-gobot/models"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

var reminders = map[time.Time][]*models.State{}
var mutexReminders sync.Mutex

// InitReminders - initialize all reminders from all states
func InitReminders() {
	for _, id := range state.ListIDs() {
		cState, initial := state.Get(id)
		if initial {
			continue
		}
		AddReminders(cState)
	}
}

// AddReminders - add all reminders from config and a single state
func AddReminders(cState *models.State) {
	deal := cState.DealOn
	if deal.Before(time.Now()) { // ignore if deal date is before today
		return
	}
	for _, offset := range config.Config.Reminders {
		addReminder(cState, offset)
	}
}

// TODO: leave review message
// addReminder - add a single reminder using state and offset
func addReminder(cState *models.State, offsetHours int8) {
	mutexReminders.Lock()
	defer mutexReminders.Unlock()

	reminderTime := cState.DealOn.Add(time.Duration(-offsetHours) * time.Hour)
	if reminderTime.Before(time.Now()) { // ignore if time after offset is before today
		return
	}

	_, exist := reminders[reminderTime]
	if exist {
		for _, id := range reminders[reminderTime] {
			if id == cState { // ignore if already exist in timeslot
				return
			}
		}
	}

	reminders[reminderTime] = append(reminders[reminderTime], cState)
	go func(stateID string) {
		if debug {
			log.Printf("Reminder logic to run for `%s` at: %s", cState.ID, reminderTime.Format("Monday, 02 January 2006, 03:04:05PM"))
		}
		<-time.After(time.Until(reminderTime))
		mutexReminders.Lock()
		defer mutexReminders.Unlock()
		for _, cState := range reminders[reminderTime] {
			until := time.Until(cState.DealOn)
			minute := int8(math.Round(until.Minutes()))
			if reminderConfigContains(config.Config.Reminders, minute) {
				hours := int8(math.Round(until.Hours()))
				message := strings.ReplaceAll(config.Config.MessageTemplates.Reminder, "{{HOURS}}", strconv.Itoa(int(hours)))

				_, err := SendMessage(cState.ID, message)
				if err != nil {
					log.Println(err)
				}

				for i, forwarder := range forwarders.Forwarders {
					message = strings.ReplaceAll(config.Config.Forwarders[i].MessageTemplates.Reminder, "{{HOURS}}", strconv.Itoa(int(hours)))
					message = strings.ReplaceAll(message, "{{ITEM}}", forwarder.Escape(cState.Name))
					message = strings.ReplaceAll(message, "{{OFFER}}", fmt.Sprintf("%.02f", cState.Price))
					err = forwarder.SendMessage(message)
					if err != nil {
						fmt.Println(err)
						continue
					}
				}
			}
		}
		delete(reminders, reminderTime)
	}(cState.ID)
}

// utilities
func reminderConfigContains(s []int8, e int8) bool {
	for _, a := range s {
		if a*60 == e {
			return true
		}
	}
	return false
}
