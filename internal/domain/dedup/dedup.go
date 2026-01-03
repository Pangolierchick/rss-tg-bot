package dedup

import (
	"context"
	"errors"

	"github.com/Pangolierchick/rss-tg-bot/internal/model"
	v2 "github.com/Pangolierchick/rss-tg-bot/pkg/rss/v2"
	"github.com/jackc/pgx/v5"
)

type repostory interface {
	GetItem(ctx context.Context, subscription int64, itemID string) (model.Item, error)
}

type Deduplicator struct {
	repo repostory
}

func New(repo repostory) *Deduplicator {
	return &Deduplicator{
		repo: repo,
	}
}

func (d *Deduplicator) IsNew(ctx context.Context, subscriptionID int64, item v2.Item) (bool, error) {
	_, err := d.repo.GetItem(ctx, subscriptionID, item.ID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return true, nil
		}

		return false, err
	}

	return false, nil
}
