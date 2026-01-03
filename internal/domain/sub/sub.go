package sub

import (
	"context"

	"github.com/Pangolierchick/rss-tg-bot/internal/model"
	v2 "github.com/Pangolierchick/rss-tg-bot/pkg/rss/v2"
)

type deduplicator interface {
	IsNew(ctx context.Context, subscriptionID int64, item v2.Item) (bool, error)
}

type notifier interface {
	Send(ctx context.Context, id int64, title string, description string, url string) error
}

type repository interface {
	AddItem(ctx context.Context, item model.Item) error

	GetSubscriptions(ctx context.Context) ([]model.Subscription, error)

	UpdatePolltime(ctx context.Context, IDs []int64) error
}

type Subscriptioner struct {
	repo     repository
	producer notifier
	rss      *v2.Fetcher
	dedup    deduplicator
}

func New(repo repository, producer notifier, rss *v2.Fetcher, dedup deduplicator) *Subscriptioner {
	return &Subscriptioner{
		repo:     repo,
		producer: producer,
		rss:      rss,
		dedup:    dedup,
	}
}
