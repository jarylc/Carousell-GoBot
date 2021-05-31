package messaging

import (
	"carousell-gobot/constants"
	"carousell-gobot/data/config"
	"carousell-gobot/models"
	"carousell-gobot/models/requests"
	"carousell-gobot/models/responses"
	"carousell-gobot/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"strconv"
	"strings"
	"time"
)

type Carousell struct {
	WS     *websocket.Conn
	ChatID string
}

func NewCarousell(ws *websocket.Conn, chatID string) Carousell {
	return Carousell{
		WS: ws,
		ChatID: chatID,
	}
}

func (m Carousell) SendMessage(text string) error {
	if text == "" {
		return nil
	}

	replyData := requests.ReplyData{
		OfferID: "",
		Source:  "web",
	}
	replyDataJSON, _ := json.Marshal(replyData)

	channelURL := strings.ReplaceAll(constants.CAROUSELL_CHANNEL, "{{CHANNEL}}", strings.ToUpper(constants.CHANNEL))
	channelURL = strings.ReplaceAll(channelURL, "{{CHATID}}", m.ChatID)
	reply := requests.Reply{
		ChannelURL:       channelURL,
		Message:          text,
		Data:             string(replyDataJSON),
		MentionType:      "users",
		MentionedUserIds: []string{},
		CustomType:       constants.MESSAGE,
		ReqID:            utils.GetEpoch(),
	}
	replyJSON, _ := json.Marshal(reply)
	replyJSONString := string(replyJSON)

	final := strings.ReplaceAll(constants.CAROUSELL_MESG, "{DATA}", replyJSONString)
	err := m.WS.WriteMessage(websocket.TextMessage, []byte(final))
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	return nil
}

func (m Carousell) Escape(str string) string {
	// no need for Carousell
	return str
}

func (m Carousell) CheckAndSendPriceMessage(info responses.MessageInfo, msg responses.Message, cState *models.State, flags *[]string, price float64) (bool, error) {
	if price == 0 {
		return false, nil
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
		return false, err
	}

	if price < listedPrice*0.5 {
		*flags = append(*flags, constants.SUPER_LOW_BALL)
	}

	if price <= listedPrice*config.Config.Carousell.LowBall {
		reply += "\n\n" + strings.ReplaceAll(config.Config.MessageTemplates.LowBalled, "{{PERCENT}}", fmt.Sprintf("%.0f", (listedPrice-price)/listedPrice*100))
		*flags = append(*flags, constants.LOW_BALL)
	}

	err = m.SendMessage(reply)
	if err != nil {
		return false, err
	}

	return true, nil
}
