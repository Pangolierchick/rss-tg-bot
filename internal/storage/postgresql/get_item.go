package postgresql

import (
	"context"

	"github.com/Pangolierchick/rss-tg-bot/internal/model"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetItem(ctx context.Context, subscriptionID int64, itemID string) (model.Item, error) {
	q := `
select
	i.item_id,
	i.subscription_id,
	i.created_at
from items i
where
	i.item_id = $1 and
	i.subscription_id = $2
	`

	rows, err := r.pool.Query(ctx, q, itemID, subscriptionID)

	if err != nil {
		return model.Item{}, err
	}

	item, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[model.Item])

	if err != nil {
		return model.Item{}, err
	}

	return *item, nil
}
