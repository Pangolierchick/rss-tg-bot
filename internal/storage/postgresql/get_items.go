package postgresql

import (
	"context"
	"time"

	"github.com/Pangolierchick/rss-tg-bot/internal/model"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetItems(ctx context.Context, subscriptionID int64, from time.Time) ([]model.Item, error) {
	q := `
select
	i.item_id,
	i.subscription_id,
	i.created_at
from items i
where
	i.subscription_id = $1 and
	i.created_at > $2
	`

	rows, err := r.pool.Query(ctx, q, subscriptionID, from)

	if err != nil {
		return nil, err
	}

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Item])

	if err != nil {
		return nil, err
	}

	return items, nil
}
