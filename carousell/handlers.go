package carousell

import (
	"carousell-gobot/constants"
	"carousell-gobot/data/config"
	"carousell-gobot/data/state"
	"carousell-gobot/messaging"
	"carousell-gobot/models/responses"
	"carousell-gobot/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

//nolint:gocognit, funlen
func handleSelling(carousellMessaging messaging.Carousell, info responses.MessageInfo, msg responses.Message, data responses.MessageData) error {
	var cState, initial = state.Get(data.OfferID)

	toForward := false
	var price float64
	var flags []string

	cState.ID = data.OfferID
	cState.Name = info.Product.Title
	cState.LastActivity = time.Now()

	switch msg.CustomType {
	case constants.MESSAGE:
		userID, err := getUserIDFromCacheOrCookie()
		if err != nil {
			return err
		}
		if msg.User.GuestID != userID { // by other party
			cState.LastReceived = msg.Message
			if initial {
				carousellMessaging.SendMessage(config.Config.MessageTemplates.FAQ)

				reply := strings.ReplaceAll(config.Config.MessageTemplates.Initial, "{{NAME}}", info.User.Username)
				reply = strings.ReplaceAll(reply, "{{ITEM}}", info.Product.Title)
				carousellMessaging.SendMessage(reply)

				toForward = true
				flags = append(flags, constants.NEW_CHAT)
			}

			switch {
			case info.IsProductSold || info.Product.Status == "R" || info.Product.Status == "D": // product sold, reserved or deleted
				reason := "sold"
				if info.Product.Status == "R" {
					reason = "reserved"
				} else if info.Product.Status == "D" {
					reason = "deleted"
				}
				carousellMessaging.SendMessage(strings.ReplaceAll(config.Config.MessageTemplates.NotAvailable, "{{REASON}}", reason))
				toForward = false
			case info.LatestPriceFormatted == "0" || info.State == "D" || info.State == "C" || debug: // official offer price $0, offer was declined or cancelled
				if info.State != "A" { // not accepted yet
					price, err = utils.GetPriceFromMessage(msg.Message)
					if err != nil {
						return err
					}
					sent, err := carousellMessaging.CheckAndSendPriceMessage(info, msg, cState, &flags, price)
					if err != nil {
						return err
					}
					if sent {
						toForward = true
					}
				}
			default:
				flags = append(flags, constants.OFFICIAL)
			}
		} else { // by myself
			cState.LastSent = msg.Message

			err = handleCommand(carousellMessaging, info, msg, data)
			if err != nil {
				return err
			}
		}
	case constants.MAKE_OFFER:
		price, err := strconv.ParseFloat(data.OfferAmount, 64)
		if err != nil {
			return err
		}
		_, err = carousellMessaging.CheckAndSendPriceMessage(info, msg, cState, &flags, price)
		if err != nil {
			return err
		}
		if price < cState.Price {
			carousellMessaging.SendMessage(config.Config.MessageTemplates.LowerOffer)
			flags = append(flags, constants.LOWERED)
		}

		toForward = true
		flags = append(flags, constants.OFFICIAL)
	}

	if toForward && !contains(flags, constants.SUPER_LOW_BALL) {
		for i, forwarder := range messaging.Forwarders {
			forward := strings.ReplaceAll(config.Config.Forwarders[i].MessageTemplates.Standard, "{{NAME}}", forwarder.Escape(msg.User.Name))
			forward = strings.ReplaceAll(forward, "{{ITEM}}", forwarder.Escape(info.Product.Title))
			forward = strings.ReplaceAll(forward, "{{ID}}", data.OfferID)
			forward = strings.ReplaceAll(forward, "{{OFFER}}", fmt.Sprintf("%.02f", price))
			forward = strings.ReplaceAll(forward, "{{FLAGS}}", strings.Join(flags, " | "))
			forwarder.SendMessage(forward)
		}
	}

	return nil
}
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func handleBuying(carousellMessaging messaging.Carousell, info responses.MessageInfo, msg responses.Message, data responses.MessageData) error {
	var cState, _ = state.Get(data.OfferID)
	cState.ID = data.OfferID
	cState.Name = info.Product.Title
	cState.LastActivity = time.Now()

	userID, err := getUserIDFromCacheOrCookie()
	if err != nil {
		return err
	}

	switch msg.CustomType {
	case constants.MESSAGE:
		cState.Price, err = utils.GetPriceFromMessage(msg.Message)
		if err != nil {
			return err
		}
		if msg.User.GuestID != userID && msg.User.Name != "" { // by other party
			cState.LastReceived = msg.Message
		} else {
			cState.LastSent = msg.Message
			err = handleCommand(carousellMessaging, info, msg, data)
			if err != nil {
				return err
			}
		}
	case constants.MAKE_OFFER:
		cState.Price, err = strconv.ParseFloat(data.OfferAmount, 64)
		if err != nil {
			return err
		}
	}
	return nil
}

func handle(raw []byte) error {
	msgString := string(raw)
	if msgString[0:4] == "MESG" {
		var msg responses.Message
		err := json.Unmarshal(raw[4:], &msg)
		if err != nil {
			return err
		}
		var data responses.MessageData
		err = json.Unmarshal([]byte(msg.Data), &data)
		time.Sleep(2 * time.Second) // wait just in case message in archive
		if err != nil {
			return err
		}
		selling, err := GetMessages(false)
		if err != nil {
			return err
		}
		carousellMessaging := messaging.NewCarousell(Connect(), data.OfferID)
		info, ok := selling[data.OfferID]
		if ok {
			err = handleSelling(carousellMessaging, info, msg, data)
			if err != nil {
				return err
			}
			err = state.Save()
			if err != nil {
				return err
			}
		}
		buying, err := GetMessages(true)
		if err != nil {
			return err
		}
		info, ok = buying[data.OfferID]
		if ok {
			err = handleBuying(carousellMessaging, info, msg, data)
			if err != nil {
				return err
			}
			err = state.Save()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
