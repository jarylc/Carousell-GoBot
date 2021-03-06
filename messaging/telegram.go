package messaging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Telegram struct {
	Token  string
	ChatID string
}

type TelegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func (m Telegram) ProcessMessage(text string) error {
	var url = "https://api.telegram.org/bot" + m.Token + "/sendMessage"

	values := TelegramMessage{
		ChatID:    m.ChatID,
		Text:      text,
		ParseMode: "Markdown",
	}
	msg, err := json.Marshal(values)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(msg)) //nolint:gosec
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
		return fmt.Errorf("telegram returned %d", resp.StatusCode)
	}

	return nil
}

func (m Telegram) Escape(str string) string {
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

func (m Telegram) SendMessage(text string) {
	addQueue(queueItem{
		messager: m,
		message:  text,
	})
}
