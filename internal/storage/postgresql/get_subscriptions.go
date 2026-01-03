package postgresql

import (
	"context"

	"github.com/Pangolierchick/rss-tg-bot/internal/model"
	"github.com/jackc/pgx/v5"
)

const (
	GetSubscriptionsLimit = 50
)

func (r *Repository) GetSubscriptions(ctx context.Context) ([]model.Subscription, error) {
	q := `
select
	s.subscription_id,
	s.user_id,
	s.url,
	s.last_polled,
	s.created_at
from subscriptions s
order by last_polled asc nulls first
limit $1
	`

	rows, err := r.pool.Query(ctx, q, GetSubscriptionsLimit)

	if err != nil {
		return nil, err
	}

	subscriptions, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Subscription])

	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}
