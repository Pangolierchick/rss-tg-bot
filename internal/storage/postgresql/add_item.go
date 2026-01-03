package postgresql

import (
	"context"

	"github.com/Pangolierchick/rss-tg-bot/internal/model"
)

func (r *Repository) AddItem(ctx context.Context, item model.Item) error {
	q := `
insert into items (item_id, subscription_id)
values ($1, $2)
	`

	_, err := r.pool.Exec(ctx, q, item.ID, item.SubscriptionID)

	return err
}
