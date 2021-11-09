package carousell

import (
	"carousell-gobot/constants"
	"carousell-gobot/data/config"
	"carousell-gobot/messaging"
	"carousell-gobot/models"
	"carousell-gobot/models/responses"
	"carousell-gobot/utils"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/dlclark/regexp2"
	"github.com/gorilla/websocket"
	"github.com/jarylc/go-chromedpproxy"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var debug = os.Getenv("DEBUG") == "1"

var ws *websocket.Conn
var mutex sync.Mutex
var mutexLocked = false

var interrupt = make(chan os.Signal, 1)

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

	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	defer func() {
		if r := recover(); r != nil {
			if userName == "" {
				userName = config.Config.Carousell.Username
			}
			if userName == "" {
				userName = "UNKNOWN"
			}
			messaging.Announce(fmt.Sprintf("`%s` errored out: %s", userName, r))
			log.Fatalln(r)
		}
	}()
main:
	for {
		token, err := getToken()
		if err != nil {
			log.Panicf("error getting token: %s", err)
		}
		userID, err := getUserIDFromCacheOrCookie()
		if err != nil {
			log.Panicf("error getting userid from cookie : %s", err)
		}
		query := strings.ReplaceAll(constants.QUERY, "{{CHANNEL}}", strings.ToUpper(constants.CHANNEL))
		query = strings.ReplaceAll(query, "{{USERID}}", userID)
		query = strings.ReplaceAll(query, "{{TOKEN}}", token)
		query = strings.ReplaceAll(query, "{{TIME}}", utils.GetEpochString())

		ws, _, err = websocket.DefaultDialer.Dial("wss://"+constants.CAROUSELL_URL_CHAT+"?"+query, nil) //nolint:bodyclose
		if err != nil {
			log.Panicf("error dialing websocket: %s", err)
		}
		if mutexLocked {
			mutexLocked = false
			mutex.Unlock()
		}

		done := make(chan struct{}, 1)
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
		log.Println("chat connected")

		pinger := time.NewTicker(time.Duration(config.Config.Carousell.PingInterval) * time.Second)
		for {
			select {
			case <-done:
				log.Println("error occurred, reconnecting in 5 seconds...")
				time.Sleep(5 * time.Second)
				continue main
			case <-pinger.C:
				err := ws.WriteMessage(websocket.TextMessage, []byte(strings.ReplaceAll(constants.CAROUSELL_PING, "{{TIME}}", utils.GetEpochString())))
				if err != nil {
					log.Println("ping error, reconnecting in 5 seconds...")
					time.Sleep(5 * time.Second)
					continue main
				}
			case <-interrupt:
				log.Println("gracefully shutting down...")
				pinger.Stop()

				err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseMessage, ""))
				if err != nil {
					log.Fatalln("error gracefully shutting down")
				}

				err = ws.Close()
				if err != nil {
					log.Fatalln("error gracefully shutting down")
				}

				select {
				case <-done:
				case <-time.After(10 * time.Second):
					log.Fatalln("graceful shutdown timed-out, forcing termination")
				case <-interrupt:
					log.Fatalln("forcefully shutting down")
				}
				log.Println("shutdown complete!")
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
var userName = ""

func getUserIDFromCacheOrCookie() (string, error) {
	if userID != "" {
		return userID, nil
	}

	r, err := regexp2.Compile("(?<=jwt=.+\\.)[^.]*", 0)
	if err != nil {
		return "", err
	}

	cUserID, err := r.FindStringMatch(config.Config.Carousell.Cookie)
	if err != nil {
		return "", err
	}

	if cUserID.String() == "" {
		return "", errors.New("unable to find user ID in cookie")
	}

	decoded, err := base64.RawStdEncoding.DecodeString(cUserID.String())
	if err != nil {
		return "", err
	}

	var jwt models.JWT
	err = json.Unmarshal(decoded, &jwt)
	if err != nil {
		return "", err
	}

	userID = jwt.ID
	userName = jwt.User
	return userID, nil
}

func getToken() (string, error) {
	var token responses.Token
	var err error

	if config.Config.Carousell.Cookie != "" {
		err = utils.HTTPGet(constants.CAROUSELL_URL_TOKEN, &token)
		if err != nil {
			return "", err
		}
	}
	if token.Data.Token == "" {
		config.Config.Carousell.Cookie, err = login()
		if err != nil {
			return "", err
		}
		err := utils.HTTPGet(constants.CAROUSELL_URL_TOKEN, &token)
		if err != nil {
			return "", err
		}
	}

	return token.Data.Token, nil
}

//nolint:funlen
func login() (string, error) {
	if config.Config.Carousell.Username == "" || config.Config.Carousell.Password == "" {
		return "", errors.New("no credentials found")
	}
	log.Print("attempting to login using credentials...")

	type result struct {
		Cookie string
		Error  error
	}
	ch := make(chan result, 1)

	chromedpproxy.PrepareProxy(config.Config.Application.ChromeListener, config.Config.Application.PortalListener, chromedp.NoSandbox)
	targetID, err := chromedpproxy.NewTab("https://www.carousell.sg/login")
	if err != nil && !errors.Is(err, context.Canceled) {
		return "", err
	}
	defer chromedpproxy.CloseTarget(targetID)
	ctx := chromedpproxy.GetTarget(targetID)

	// input credentials
	err = chromedp.Run(ctx, chromedp.Tasks{
		chromedp.WaitReady(`.grecaptcha-badge`),
		chromedp.SendKeys(`input[name="username"]`, config.Config.Carousell.Username, chromedp.NodeVisible),
		chromedp.SendKeys(`input[name="password"]`, config.Config.Carousell.Password, chromedp.NodeVisible),
		chromedp.Click(`button[type="submit"]`, chromedp.NodeVisible),
	})
	if err != nil && !errors.Is(err, context.Canceled) {
		return "", err
	}

	// wrong username/password
	go func() {
		err := chromedp.Run(ctx, chromedp.Tasks{
			chromedp.WaitVisible(`//p[contains(text(), 'Wrong username or password')]/..`, chromedp.NodeVisible),
			chromedp.ActionFunc(func(ctx context.Context) error {
				return errors.New("invalid credentials")
			}),
		})
		if err != nil && !errors.Is(err, context.Canceled) {
			ch <- result{
				Cookie: "",
				Error:  err,
			}
		}
	}()

	// 2FA is required
	go func() {
		err := chromedp.Run(ctx, chromedp.Tasks{
			chromedp.WaitVisible(`input[name="verification code"]`),
			chromedp.ActionFunc(func(ctx context.Context) error {
				msg := fmt.Sprintf("2FA required, please solve it: %s/?id=%s", config.Config.Application.BaseURL, targetID)
				log.Print(msg)
				messaging.Announce(msg)
				return nil
			}),
		})
		if err != nil && !errors.Is(err, context.Canceled) {
			ch <- result{
				Cookie: "",
				Error:  err,
			}
		}
	}()

	// success
	go func() {
		err := chromedp.Run(ctx, chromedp.Tasks{
			chromedp.WaitVisible(`//p[contains(text(), 'Hello,')]/..`, chromedp.NodeVisible),
			chromedp.ActionFunc(func(ctx context.Context) error {
				cookies, err := network.GetCookies().WithUrls([]string{"https://www.carousell.sg", "https://carousell.sg"}).Do(ctx)
				if err != nil {
					return err
				}
				cookieStr := ""
				for _, cookie := range cookies {
					cookieStr += fmt.Sprintf("%s=%s; ", cookie.Name, cookie.Value)
				}
				ch <- result{
					Cookie: cookieStr,
					Error:  nil,
				}
				return nil
			}),
		})
		if err != nil && !errors.Is(err, context.Canceled) {
			ch <- result{
				Cookie: "",
				Error:  err,
			}
		}
	}()

	final := <-ch
	return final.Cookie, final.Error
}
