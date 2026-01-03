package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type sendMessageBody struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

type Telegram struct {
	http     *http.Client
	botToken string
	chatID   int64
}

func NewTelegam(http *http.Client, botToken string) *Telegram {
	return &Telegram{
		http:     http,
		botToken: botToken,
	}
}

func (t *Telegram) Send(ctx context.Context, id int64, title string, description string, url string) error {
	b := sendMessageBody{
		ChatID: id,
		Text:   fmt.Sprintf("%s\n\n%s", description, url),
	}
	json, err := json.Marshal(b)

	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.botToken), bytes.NewBuffer(json))
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return err
	}

	resp, err := t.http.Do(req)
	if err != nil {
		return err
	}

	slog.Debug("telegram response",
		"status", resp.StatusCode,
		"body", resp.Body,
	)

	defer resp.Body.Close()

	return nil
}
