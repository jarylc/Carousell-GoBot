package forwarders

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Slack struct {
	WebhookURL string
}

type SlackMessage struct {
	Text string `json:"text"`
}

func (m Slack) SendMessage(text string) error {
	values := SlackMessage{
		Text: text,
	}
	msg, err := json.Marshal(values)
	if err != nil {
		return err
	}

	resp, err := http.Post(m.WebhookURL, "application/json", bytes.NewReader(msg))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		log.Println(string(body))
		return errors.New("slack returned non-200 status code")
	}

	return nil
}

func (m Slack) Escape(str string) string {
	// mrkdwn
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"[", "\\[",
		"<", "&lt;",
		">", "&gt;",
	)
	str = replacer.Replace(str)
	return str
}
