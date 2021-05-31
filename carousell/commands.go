package carousell

import (
	"carousell-gobot/chrono"
	"carousell-gobot/data/config"
	"carousell-gobot/data/state"
	"carousell-gobot/messaging"
	"carousell-gobot/models/responses"
	"errors"
	"fmt"
	"github.com/dlclark/regexp2"
	"log"
	"strings"
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
	case "sched", "schedule", "remind", "reminder", "deal": // schedule
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
			AddReminders(cState)
			messaging.SendMessage(fmt.Sprintf("Deal scheduled on: %s\nReminders set: %shr(s) before", parse.Format("Monday, 02 January 2006, 03:04:05PM"), strings.Trim(strings.Join(strings.Fields(fmt.Sprint(config.Config.Reminders)), "hr(s), "), "[]")))
		} else {
			return errors.New("ERROR: Unable to parse date from messages")
		}
	case "faq": // resend faq
		messaging.SendMessage(config.Config.MessageTemplates.FAQ)
	}

	return nil
}
