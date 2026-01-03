package tgsender

import (
	"context"
	"fmt"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"golang.org/x/time/rate"
)

type TgSender struct {
	bot     *bot.Bot
	limiter *rate.Limiter
}

func New(bot *bot.Bot) *TgSender {
	limiter := rate.NewLimiter(rate.Every(1*time.Second), 30)
	return &TgSender{
		bot:     bot,
		limiter: limiter,
	}
}

func (t *TgSender) Send(ctx context.Context, ID any, message string) error {
	if t.limiter.Allow() == false {
		return fmt.Errorf("failed to send tg message: rate limit exceded")
	}
	_, err := t.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    ID,
		Text:      message,
		ParseMode: models.ParseModeMarkdownV1,
	})

	return err
}
