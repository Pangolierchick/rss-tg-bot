package repository

import (
	"context"

	"github.com/Pangolierchick/rss-tg-bot/internal/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetSubscriber(ctx context.Context, subscriberID int64) (*models.Subscriber, error) {
	q := `
select
	subscriber_id,
	tg_chat_id,
	created_at
from subscribers
where
	subscriber_id = $1
	`

	rows, err := r.pool.Query(ctx, q, subscriberID)

	if err != nil {
		return nil, err
	}

	subscriber, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[models.Subscriber])

	return subscriber, err
}
