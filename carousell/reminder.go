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

var reminders = map[*models.State][]*models.Reminder{}
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

	// cancel existing reminders
	cReminders, exist := reminders[cState]
	if exist {
		for _, reminders := range cReminders {
			reminders.Cancel()
		}
	}

	for _, offset := range config.Config.Reminders {
		addReminder(cState, offset)
	}
}

// TODO: leave review message
// addReminder - add a single reminder using state and offset
func addReminder(cState *models.State, offsetHours int16) {
	mutexReminders.Lock()
	defer mutexReminders.Unlock()

	reminderTime := cState.DealOn.Add(time.Duration(-offsetHours) * time.Hour)
	if reminderTime.Before(time.Now()) { // ignore if time after offset is before today
		return
	}

	reminder := models.NewReminder(reminderTime)
	reminders[cState] = append(reminders[cState], reminder)

	go func(cState *models.State, reminder *models.Reminder) {
		if debug {
			log.Printf("Reminder to run for `%s` at: %s", cState.ID, reminder.Time.Format("Monday, 02 January 2006, 03:04:05PM"))
		}
		select {
		case <-time.After(time.Until(reminder.Time)):
			if debug {
				log.Printf("Reminder ran for `%s`", cState.ID)
			}

			reminder.Close()

			until := time.Until(cState.DealOn)
			hours := int16(math.Round(until.Hours()))

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
		case <-reminder.ChanCancel:
			if debug {
				log.Printf("Reminder cancelled for `%s` at: %s", cState.ID, reminder.Time.Format("Monday, 02 January 2006, 03:04:05PM"))
			}
		}
		mutexReminders.Lock()
		delete(reminders, cState)
		mutexReminders.Unlock()
	}(cState, reminder)
}
