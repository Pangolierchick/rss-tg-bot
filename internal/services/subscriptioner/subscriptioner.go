package subscriptioner

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Pangolierchick/rss-tg-bot/internal/models"
)

var (
	SubscriptionsLimit int64 = 20
)

type repository interface {
	CountSubscriptions(ctx context.Context, tgChatId int64) (int64, error)

	AddSubscription(ctx context.Context, subscription *models.Subscription) error
	AddSubscriber(ctx context.Context, subscriber *models.Subscriber) (int64, error)
	AddFeed(ctx context.Context, feed *models.Feed) (int64, error)

	GetFeedByURL(ctx context.Context, URL string) (*models.Feed, error)
	GetSubscriberByTg(ctx context.Context, tgChatID int64) (*models.Subscriber, error)
	GetUserSubscriptionsURLs(ctx context.Context, tgChatID int64) ([]string, error)

	DeleteSubscription(ctx context.Context, feedID, subscriberID int64) error
}

type Subscriptioner struct {
	repo repository
}

func New(repo repository) *Subscriptioner {
	return &Subscriptioner{
		repo: repo,
	}
}

type AddSubscriptionParams struct {
	TgChatID int64
	URL      string
}

func (s *Subscriptioner) AddSubscription(ctx context.Context, params *AddSubscriptionParams) error {
	count, err := s.repo.CountSubscriptions(ctx, params.TgChatID)

	if err != nil {
		slog.Error("failed to count user's subscriptions",
			"error", err,
		)
		return err
	}

	if count >= SubscriptionsLimit {
		return fmt.Errorf("user has reached subscriptions limit of %d", SubscriptionsLimit)
	}

	feedID, err := s.repo.AddFeed(ctx, &models.Feed{
		URL: params.URL,
	})

	if err != nil {
		slog.Error("failed to add feed",
			"url", params.URL,
			"error", err,
		)
		return err
	}

	subscriberID, err := s.repo.AddSubscriber(ctx, &models.Subscriber{
		TgChatID: params.TgChatID,
	})

	if err != nil {
		slog.Error("failed to add subscriber",
			"error", err,
		)
		return err
	}

	err = s.repo.AddSubscription(ctx, &models.Subscription{
		FeedID:       feedID,
		SubscriberID: subscriberID,
	})

	if err != nil {
		slog.Error("failed to add subscription",
			"feed_id", feedID,
			"subscriber_id", subscriberID,
			"error", err,
		)
	}

	return err
}

type DeleteSubscriptionParams struct {
	TgChatID int64
	URL      string
}

func (s *Subscriptioner) DeleteSubscription(ctx context.Context, params *DeleteSubscriptionParams) error {
	subscriber, err := s.repo.GetSubscriberByTg(ctx, params.TgChatID)

	if err != nil {
		slog.Error("failed to get subscriber by tg",
			"error", err,
			"tg_chat_id", params.TgChatID,
		)
		return err
	}

	feed, err := s.repo.GetFeedByURL(ctx, params.URL)

	if err != nil {
		slog.Error("failed to get feed by url",
			"error", err,
			"url", params.URL,
		)
		return err
	}

	err = s.repo.DeleteSubscription(ctx, feed.ID, subscriber.SubscriberID)

	if err != nil {
		slog.Error("failed to delete subscription",
			"error", err,
			"feed_id", feed.ID,
			"subscriber_id", subscriber.SubscriberID,
		)
	}

	return err
}

func (s *Subscriptioner) GetSubscriptions(ctx context.Context, tgChatID int64) ([]string, error) {
	URLs, err := s.repo.GetUserSubscriptionsURLs(ctx, tgChatID)

	if err != nil {
		slog.Error("failed to get users subscriptions URLs",
			"error", err,
		)
		return nil, err
	}

	return URLs, nil
}
