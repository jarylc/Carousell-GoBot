package carousell

import (
	"carousell-gobot/constants"
	"carousell-gobot/data/config"
	"carousell-gobot/data/state"
	"carousell-gobot/messaging"
	"carousell-gobot/models/responses"
	"errors"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/jarylc/go-chrono/v2"
	"log"
	"strings"
	"syscall"
	"time"
)

//nolint:funlen, gocognit
func handleCommand(messaging messaging.Carousell, info responses.MessageInfo, msg responses.Message, data responses.MessageData) error {
	var err error

	if !strings.HasPrefix(msg.Message, config.Config.CommandPrefix) {
		return nil
	} // ignore if not command

	cState, initial := state.Get(data.OfferID)
	if initial {
		return errors.New("command could not find state")
	}

	cmd := strings.TrimSpace(strings.Fields(msg.Message[1:])[0])

	regex, err := regexp2.Compile("(?<=.+ ).+", 0)
	if err != nil {
		return err // probably will not happen
	}
	args, err := regex.FindStringMatch(msg.Message)
	if err != nil {
		return err // probably will not happen
	}

	if debug {
		log.Printf("Command recieved `%s`, arguments: %s\n", cmd, args)
	}

	switch cmd {
	case "sched", "schedule", "remind", "reminder", "deal": // confirm deal
		var c chrono.Chrono
		c, err = chrono.New()
		if err != nil {
			return err
		}

		var parse *time.Time
		if args != nil { // with argument
			parse, err = c.ParseDate(args.String(), time.Now())
			if err != nil || parse == nil {
				messaging.SendMessage("ERROR: Invalid natural date")
			}
		} else {
			parse, err = c.ParseDate(cState.LastReceived, time.Now())
			if err != nil || parse == nil {
				parse, err = c.ParseDate(cState.LastSent, time.Now())
				if err != nil || parse == nil {
					messaging.SendMessage("ERROR: Unable to find natural date in last response and reply, please specify in argument")
				}
			}
		}

		if parse != nil {
			cState.DealOn = time.Unix(parse.Unix(), 0)
			reminderTimes := AddReminders(cState)
			remindersTimesStr := ""
			if len(reminderTimes) > 0 {
				remindersTimesStr = "\n\nReminders will be sent on:"
				for i, reminderTime := range reminderTimes {
					remindersTimesStr += fmt.Sprintf("\n%d. %s", i+1, reminderTime.Format(constants.READABLE_DATE_FORMAT))
				}
			}
			messaging.SendMessage(fmt.Sprintf("Deal scheduled on: %s%s", parse.Format(constants.READABLE_DATE_FORMAT), remindersTimesStr))
		} else {
			return errors.New("ERROR: Unable to parse date from messages")
		}
	case "cancel", "del", "delete": // delete deal
		cState.DealOn = time.Time{}
		CancelReminders(cState)
		messaging.SendMessage("Deal cancelled, reminders unscheduled.")
	case "faq": // resend faq
		messaging.SendMessage(config.Config.MessageTemplates.FAQ)
	case "contact": // send contact details
		messaging.SendMessage(config.Config.MessageTemplates.Contact)
	case "stop": // stop bot
		interrupt <- syscall.SIGINT
	}

	return nil
}
