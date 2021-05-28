package carousell

import (
	"carousell-gobot/constants"
	"carousell-gobot/data/config"
	"carousell-gobot/models/requests"
	"carousell-gobot/models/responses"
	"carousell-gobot/utils"
	"encoding/json"
	"github.com/dlclark/regexp2"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"
)

var debug = os.Getenv("DEBUG") == "1"

var ws *websocket.Conn
var mutex sync.Mutex

// Connect - return websocket connection, if not create it
//nolint:funlen,gocognit
func Connect() *websocket.Conn {
	mutex.Lock() // if already attempting to connect, wait
	if ws != nil {
		return ws
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	userID, err := getUserIDFromCacheOrCookie()
	if err != nil {
		log.Fatalln(err)
	}
main:
	for {
		token, err := getToken()
		if err != nil {
			log.Fatalln(err)
		}
		query := strings.ReplaceAll(constants.QUERY, "{{CHANNEL}}", strings.ToUpper(constants.CHANNEL))
		query = strings.ReplaceAll(query, "{{USERID}}", userID)
		query = strings.ReplaceAll(query, "{{TOKEN}}", token)
		query = strings.ReplaceAll(query, "{{TIME}}", utils.GetEpochString())

		ws, _, err = websocket.DefaultDialer.Dial("wss://"+constants.CAROUSELL_URL_CHAT+"?"+query, nil) //nolint:bodyclose
		if err != nil {
			log.Panic(err)
		}
		mutex.Unlock()

		done := make(chan struct{})
		go func() {
			defer close(done)
			for {
				_, message, err := ws.ReadMessage()
				if err != nil {
					break
				}
				err = handle(message)
				if err != nil {
					log.Println(err)
					break
				}
			}
		}()
		log.Println("Chat connected")

		log.Println("Initiating reminders system")
		InitReminders()
		log.Println("Reminder system initiated")

		pinger := time.NewTicker(time.Duration(config.Config.Carousell.PingInterval) * time.Second)
		for {
			select {
			case <-done:
				return ws
			case <-pinger.C:
				err := ws.WriteMessage(websocket.TextMessage, []byte(strings.ReplaceAll(constants.CAROUSELL_PING, "{{TIME}}", utils.GetEpochString())))
				if err != nil {
					log.Println("Ping error, restarting in 5 seconds...")
					time.Sleep(5 * time.Second)
					continue main
				}
			case <-interrupt:
				log.Println("Gracefully shutting down...")
				pinger.Stop()

				err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseMessage, ""))
				if err != nil {
					log.Panic("Error gracefully shutting down")
				}

				err = ws.Close()
				if err != nil {
					log.Panic("Error gracefully shutting down")
				}

				select {
				case <-done:
				case <-time.After(10 * time.Second):
					log.Fatalln("Graceful shutdown timed-out, forcing termination")
				case <-interrupt:
					log.Fatalln("Forcefully shutting down")
				}
				log.Println("Shutdown complete!")
				return ws
			}
		}
	}
}

func SendMessage(chatID string, msg string) (string, error) {
	if msg == "" {
		return "", nil
	}

	replyData := requests.ReplyData{
		OfferID: "",
		Source:  "web",
	}
	replyDataJSON, _ := json.Marshal(replyData)

	channelURL := strings.ReplaceAll(constants.CAROUSELL_CHANNEL, "{{CHANNEL}}", strings.ToUpper(constants.CHANNEL))
	channelURL = strings.ReplaceAll(channelURL, "{{CHATID}}", chatID)
	reply := requests.Reply{
		ChannelURL:       channelURL,
		Message:          msg,
		Data:             string(replyDataJSON),
		MentionType:      "users",
		MentionedUserIds: []string{},
		CustomType:       constants.MESSAGE,
		ReqID:            utils.GetEpoch(),
	}
	replyJSON, _ := json.Marshal(reply)
	replyJSONString := string(replyJSON)

	final := strings.ReplaceAll(constants.CAROUSELL_MESG, "{DATA}", replyJSONString)
	err := ws.WriteMessage(websocket.TextMessage, []byte(final))
	if err != nil {
		return "", err
	}

	time.Sleep(1 * time.Second)

	return final, nil
}

// GetMessages - get all Carousell messages from API
func GetMessages(buying bool) (map[string]responses.MessageInfo, error) {
	offers := map[string]responses.MessageInfo{}

	var messageType = "received"
	if buying {
		messageType = "made"
	}

	var messages responses.Messages
	err := utils.HTTPGet(strings.ReplaceAll(constants.CAROUSELL_URL_MESSAGES, "{{TYPE}}", messageType), &messages)
	if err != nil {
		return offers, err
	}

	for _, offer := range messages.Data.Offers {
		offers[strconv.Itoa(offer.ID)] = offer
	}
	return offers, nil
}

var userID = ""
func getUserIDFromCacheOrCookie() (string, error) {
	if userID != "" {
		return userID, nil
	}

	r, err := regexp2.Compile("(?<=_t=.*%3D)\\d+(?=;)", 0)
	if err != nil {
		return "", err
	}

	cUserID, err := r.FindStringMatch(config.Config.Carousell.Cookie)
	if err != nil {
		return "", err
	}

	userID = cUserID.String()

	return userID, nil
}

func getToken() (string, error) {
	var token responses.Token

	err := utils.HTTPGet(constants.CAROUSELL_URL_TOKEN, &token)
	if err != nil {
		return "", err
	}

	return token.Data.Token, nil
}
