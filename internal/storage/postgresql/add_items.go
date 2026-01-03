package postgresql

import (
	"context"

	"github.com/Pangolierchick/rss-tg-bot/internal/model"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) AddItems(ctx context.Context, items []model.Item) error {
	columns := []string{"item_id", "subscription_id"}

	source := pgx.CopyFromSlice(len(items), func(i int) ([]any, error) {
		return []any{
			items[i].ID,
			items[i].SubscriptionID,
		}, nil
	})

	_, err := r.pool.CopyFrom(ctx, pgx.Identifier{"items"}, columns, source)

	return err
}
