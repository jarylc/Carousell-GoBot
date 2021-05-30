package carousell

import (
	"carousell-gobot/constants"
	"carousell-gobot/data/config"
	"carousell-gobot/data/state"
	"carousell-gobot/forwarders"
	"carousell-gobot/models"
	"carousell-gobot/models/responses"
	"carousell-gobot/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

//nolint:gocognit, funlen
func handleSelling(info responses.MessageInfo, msg responses.Message, data responses.MessageData) error {
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
			cState.LastResponse = msg.Message

			if initial {
				_, err = SendMessage(data.OfferID, config.Config.MessageTemplates.FAQ)
				if err != nil {
					return err
				}

				reply := strings.ReplaceAll(config.Config.MessageTemplates.Initial, "{{NAME}}", info.User.Username)
				reply = strings.ReplaceAll(reply, "{{ITEM}}", info.Product.Title)
				_, err = SendMessage(data.OfferID, reply)
				if err != nil {
					return err
				}

				toForward = true
				flags = append(flags, constants.NEW_OFFER)
			}

			if info.LatestPriceFormatted == "0" || info.State == "D" || info.State == "C" || debug { // if official offer not made yet, declined, cancelled or debug mode
				if info.State != "A" { // not accepted yet
					price, err = utils.GetPriceFromMessage(msg.Message)
					if err != nil {
						return err
					}
					sent, err := checkAndSendPriceMessage(info, msg, data, cState, &flags, price)
					if err != nil {
						return err
					}
					if sent != "" {
						toForward = true
					}
				}
			} else {
				flags = append(flags, constants.OFFICIAL)
			}
		} else { // by myself
			cState.LastReply = msg.Message

			err = handleCommand(info, msg, data)
			if err != nil {
				return err
			}
		}
	case constants.MAKE_OFFER:
		price, err := strconv.ParseFloat(data.OfferAmount, 64)
		if err != nil {
			return err
		}
		sent, err := checkAndSendPriceMessage(info, msg, data, cState, &flags, price)
		if err != nil {
			return err
		}
		if price < cState.Price {
			_, err := SendMessage(data.OfferID, config.Config.MessageTemplates.LowerOffer)
			if err != nil {
				return err
			}
			flags = append(flags, constants.LOWERED)
		}
		if sent != "" {
			toForward = true
			flags = append(flags, constants.OFFICIAL)
		}
	}

	if toForward && !contains(flags, constants.SUPER_LOW_BALL) {
		for i, forwarder := range forwarders.Forwarders {
			forward := strings.ReplaceAll(config.Config.Forwarders[i].MessageTemplates.Standard, "{{NAME}}", forwarder.Escape(msg.User.Name))
			forward = strings.ReplaceAll(forward, "{{ITEM}}", forwarder.Escape(info.Product.Title))
			forward = strings.ReplaceAll(forward, "{{ID}}", data.OfferID)
			forward = strings.ReplaceAll(forward, "{{OFFER}}", fmt.Sprintf("%.02f", price))
			forward = strings.ReplaceAll(forward, "{{FLAGS}}", strings.Join(flags, " | "))
			err := forwarder.SendMessage(forward)
			if err != nil {
				fmt.Println(err)
				continue
			}
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

func handleBuying(info responses.MessageInfo, msg responses.Message, data responses.MessageData) error {
	var cState, _ = state.Get(data.OfferID)
	cState.ID = data.OfferID
	cState.Name = info.Product.Title
	cState.LastActivity = time.Now()

	userID, err := getUserIDFromCacheOrCookie()
	if err != nil {
		return err
	}
	if msg.User.GuestID != userID && msg.User.Name != "" { // by other party
		cState.LastResponse = msg.Message
	} else {
		cState.LastReply = msg.Message
	}

	err = handleCommand(info, msg, data)
	if err != nil {
		return err
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
		info, ok := selling[data.OfferID]
		if ok {
			err = handleSelling(info, msg, data)
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
			err = handleBuying(info, msg, data)
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

func checkAndSendPriceMessage(info responses.MessageInfo, msg responses.Message, data responses.MessageData, cState *models.State, flags *[]string, price float64) (string, error) {
	if price == -1 {
		return "", nil
	}
	cState.Price = price

	reply := ""
	if msg.CustomType == constants.MAKE_OFFER {
		reply = config.Config.MessageTemplates.Offered
	} else {
		reply = config.Config.MessageTemplates.PossibleOffer
	}

	reply = strings.ReplaceAll(reply, "{{OFFER}}", fmt.Sprintf("%.02f", price))
	listedPrice, err := strconv.ParseFloat(info.Product.Price, 64)
	if err != nil {
		return "", err
	}

	if price < listedPrice*0.5 && !debug { // don't bother treating these as potential offers, even Carousell doesn't
		*flags = append(*flags, constants.SUPER_LOW_BALL)
	}

	if price <= listedPrice*config.Config.Carousell.LowBall {
		reply += "\n\n" + strings.ReplaceAll(config.Config.MessageTemplates.LowBalled, "{{PERCENT}}", fmt.Sprintf("%.0f", (listedPrice-price)/listedPrice*100))
		*flags = append(*flags, constants.LOW_BALL)
	}

	send, err := SendMessage(data.OfferID, reply)
	if err != nil {
		return "", err
	}

	return send, nil
}
