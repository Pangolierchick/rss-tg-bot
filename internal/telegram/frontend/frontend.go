package telegramfrontend

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	s "github.com/Pangolierchick/rss-tg-bot/internal/services/subscriptioner"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	addSubscriptionCommand     = "add"
	getSubscriptionsCommand    = "get"
	deleteSubscriptionsCommand = "del"
)

type subscriptioner interface {
	GetSubscriptions(ctx context.Context, tgChatID int64) ([]string, error)
	DeleteSubscription(ctx context.Context, params *s.DeleteSubscriptionParams) error
	AddSubscription(ctx context.Context, params *s.AddSubscriptionParams) error
}

type TgFrontend struct {
	bot *bot.Bot
	sub subscriptioner
}

func New(bot *bot.Bot, sub subscriptioner) *TgFrontend {
	tg := &TgFrontend{
		bot: bot,
		sub: sub,
	}

	tg.registerHandlers()

	return tg
}

func (t *TgFrontend) Run(ctx context.Context) {
	go t.bot.Start(ctx)
}

func (t *TgFrontend) registerHandlers() {
	t.bot.RegisterHandler(bot.HandlerTypeMessageText, addSubscriptionCommand, bot.MatchTypeCommand, t.handleAddSubscription)
	t.bot.RegisterHandler(bot.HandlerTypeMessageText, getSubscriptionsCommand, bot.MatchTypeCommand, t.handleGetSubscriptions)
	t.bot.RegisterHandler(bot.HandlerTypeMessageText, deleteSubscriptionsCommand, bot.MatchTypeCommand, t.handleDeleteSubscription)
}

func (t *TgFrontend) handleAddSubscription(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	text := strings.TrimSpace(update.Message.Text)

	args := strings.Split(text, " ")
	if len(args) < 2 {
		t.sendReply(chatID, "Usage: /add <RSS_URL>")
		return
	}

	rssURL := strings.TrimSpace(args[1])
	if rssURL == "" {
		t.sendReply(chatID, "RSS URL cannot be empty")
		return
	}

	if !strings.HasPrefix(rssURL, "http://") && !strings.HasPrefix(rssURL, "https://") {
		t.sendReply(chatID, "Invalid URL format. URL must start with http:// or https://")
		return
	}

	params := &s.AddSubscriptionParams{
		TgChatID: chatID,
		URL:      rssURL,
	}

	err := t.sub.AddSubscription(ctx, params)
	if err != nil {
		t.sendReply(chatID, fmt.Sprintf("Error adding subscription: %v", err))
		return
	}

	t.sendReply(chatID, fmt.Sprintf("Successfully added subscription for: %s", rssURL))
}

func (t *TgFrontend) handleGetSubscriptions(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID

	subscriptions, err := t.sub.GetSubscriptions(ctx, chatID)
	if err != nil {
		t.sendReply(chatID, fmt.Sprintf("Error getting subscriptions: %v", err))
		return
	}

	if len(subscriptions) == 0 {
		t.sendReply(chatID, "You have no subscriptions.")
		return
	}

	var response strings.Builder
	for i, url := range subscriptions {
		response.WriteString(fmt.Sprintf("%d.\t%s\n", i+1, url))
	}

	t.sendReply(chatID, response.String())
}

func (t *TgFrontend) handleDeleteSubscription(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	text := strings.TrimSpace(update.Message.Text)

	args := strings.Split(text, " ")
	if len(args) < 2 {
		t.sendReply(chatID, "Usage: /del <RSS_URL> or /del <number>")
		return
	}

	identifier := strings.TrimSpace(args[1])
	if identifier == "" {
		t.sendReply(chatID, "Identifier cannot be empty")
		return
	}

	if idx, err := strconv.Atoi(identifier); err == nil {
		subscriptions, err := t.sub.GetSubscriptions(ctx, chatID)
		if err != nil {
			t.sendReply(chatID, fmt.Sprintf("Error getting subscriptions: %v", err))
			return
		}

		if idx < 1 || idx > len(subscriptions) {
			t.sendReply(chatID, fmt.Sprintf("Invalid subscription number. Please provide a number between 1 and %d", len(subscriptions)))
			return
		}

		urlToDelete := subscriptions[idx-1]
		params := &s.DeleteSubscriptionParams{
			TgChatID: chatID,
			URL:      urlToDelete,
		}

		err = t.sub.DeleteSubscription(ctx, params)
		if err != nil {
			t.sendReply(chatID, fmt.Sprintf("Error deleting subscription: %v", err))
			return
		}

		t.sendReply(chatID, fmt.Sprintf("Successfully deleted subscription: %s", urlToDelete))
	} else {
		if !strings.HasPrefix(identifier, "http://") && !strings.HasPrefix(identifier, "https://") {
			t.sendReply(chatID, "Invalid URL format. URL must start with http:// or https:// or provide a valid subscription number")
			return
		}

		params := &s.DeleteSubscriptionParams{
			TgChatID: chatID,
			URL:      identifier,
		}

		err := t.sub.DeleteSubscription(ctx, params)
		if err != nil {
			t.sendReply(chatID, fmt.Sprintf("Error deleting subscription: %v", err))
			return
		}

		t.sendReply(chatID, fmt.Sprintf("Successfully deleted subscription: %s", identifier))
	}
}

func (t *TgFrontend) sendReply(chatID int64, text string) {
	_, err := t.bot.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		slog.Error("Error sending message to chat", "chat_id", chatID, "error", err)
	}
}
