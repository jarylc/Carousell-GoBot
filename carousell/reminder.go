package carousell

import (
	"carousell-gobot/constants"
	"carousell-gobot/data/config"
	"carousell-gobot/data/state"
	"carousell-gobot/messaging"
	"carousell-gobot/models"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

var reminders = map[string][]*models.Reminder{} // id as key
var mutexReminders sync.Mutex

// InitReminders - initialize all reminders from all states
func InitReminders() {
	for _, id := range state.ListIDs() {
		cState, initial := state.Get(id)
		if initial {
			continue
		}
		_ = AddReminders(cState)
	}
}

// AddReminders - add all reminders from config and a single state
func AddReminders(cState *models.State) []time.Time {
	deal := cState.DealOn
	if deal.Before(time.Now()) { // ignore if deal date is before today
		return nil
	}

	// cancel existing reminders
	CancelReminders(cState)

	var reminderTimes []time.Time
	for _, offset := range config.Config.Reminders {
		reminderTime := addReminder(cState, offset)
		if !reminderTime.IsZero() {
			reminderTimes = append(reminderTimes, reminderTime)
		}
	}
	return reminderTimes
}

// CancelReminders - cancel all reminders for a state
func CancelReminders(cState *models.State) {
	cReminders, exist := reminders[cState.ID]
	if exist {
		for _, reminders := range cReminders {
			reminders.Cancel()
		}
	}
}

// TODO: leave review message
// addReminder - add a single reminder using state and offset
func addReminder(cState *models.State, offsetHours int16) time.Time {
	mutexReminders.Lock()
	defer mutexReminders.Unlock()

	reminderTime := cState.DealOn.Add(time.Duration(-offsetHours) * time.Hour)
	if reminderTime.Before(time.Now()) { // ignore if time after offset is before today
		return time.Time{}
	}

	reminder := models.NewReminder(reminderTime)
	reminders[cState.ID] = append(reminders[cState.ID], reminder)

	go func(cState *models.State, reminder *models.Reminder) {
		if debug {
			log.Printf("Reminder to run for `%s` at: %s", cState.ID, reminder.Time.Format(constants.READABLE_DATE_FORMAT))
		}

		delay := time.NewTimer(time.Until(reminder.Time))
		select {
		case <-delay.C:
			if debug {
				log.Printf("Reminder ran for `%s`", cState.ID)
			}

			reminder.Close()

			until := time.Until(cState.DealOn)
			hours := int16(math.Round(until.Hours()))

			if hours > 0 {
				message := strings.ReplaceAll(config.Config.MessageTemplates.Reminder, "{{HOURS}}", strconv.Itoa(int(hours)))
				messaging.NewCarousell(Connect(), cState.ID).SendMessage(message)

				for i, forwarder := range messaging.Forwarders {
					message = strings.ReplaceAll(config.Config.Forwarders[i].MessageTemplates.Reminder, "{{HOURS}}", strconv.Itoa(int(hours)))
					message = strings.ReplaceAll(message, "{{ITEM}}", forwarder.Escape(cState.Name))
					message = strings.ReplaceAll(message, "{{ID}}", forwarder.Escape(cState.ID))
					message = strings.ReplaceAll(message, "{{OFFER}}", fmt.Sprintf("%.02f", cState.Price))
					forwarder.SendMessage(message)
				}
			}
		case <-reminder.ChanCancel:
			if !delay.Stop() {
				<-delay.C
			}
			if debug {
				log.Printf("Reminder cancelled for `%s` at: %s", cState.ID, reminder.Time.Format(constants.READABLE_DATE_FORMAT))
			}
		}
		mutexReminders.Lock()
		delete(reminders, cState.ID)
		mutexReminders.Unlock()
	}(cState, reminder)

	return reminderTime
}
