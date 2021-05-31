package carousell

import (
	"carousell-gobot/constants"
	"carousell-gobot/data/config"
	"carousell-gobot/models/responses"
	"carousell-gobot/utils"
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
var mutexLocked = false

// Connect - return websocket connection, if not create it
//nolint:funlen,gocognit

func Connect() *websocket.Conn {
	mutex.Lock()
	mutexLocked = true
	if ws != nil {
		mutexLocked = false
		mutex.Unlock()
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
		if mutexLocked {
			mutexLocked = false
			mutex.Unlock()
		}

		done := make(chan struct{})
		go func() {
			defer close(done)
			for {
				_, message, err := ws.ReadMessage()
				if err != nil {
					log.Println(err)
					return
				}
				err = handle(message)
				if err != nil {
					log.Println(err)
					return
				}
			}
		}()
		log.Println("Chat connected")

		pinger := time.NewTicker(time.Duration(config.Config.Carousell.PingInterval) * time.Second)
		for {
			select {
			case <-done:
				log.Println("Error occurred, reconnecting in 5 seconds...")
				time.Sleep(5 * time.Second)
				continue main
			case <-pinger.C:
				err := ws.WriteMessage(websocket.TextMessage, []byte(strings.ReplaceAll(constants.CAROUSELL_PING, "{{TIME}}", utils.GetEpochString())))
				if err != nil {
					log.Println("Ping error, reconnecting in 5 seconds...")
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

	r, err := regexp2.Compile("(?<=_t=.*(u%3D|u=))\\d+(?=;)", 0)
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
