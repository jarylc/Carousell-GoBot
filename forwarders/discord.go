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

type Discord struct {
	WebhookURL string
}

type DiscordMessage struct {
	Content string `json:"content"`
}

func (m Discord) SendMessage(text string) error {
	values := DiscordMessage{
		Content: text,
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
		return errors.New("discord returned non-200 status code")
	}

	return nil
}

func (m Discord) Escape(str string) string {
	// Markdown
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"`", "\\`",
		"\\", "\\\\",
	)
	str = replacer.Replace(str)
	return str
}