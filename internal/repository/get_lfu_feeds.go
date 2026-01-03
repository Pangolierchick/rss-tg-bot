package repository

import (
	"context"

	"github.com/Pangolierchick/rss-tg-bot/internal/models"
	"github.com/jackc/pgx/v5"
)

type GetLFUFeedsParams struct {
	Limit int
}

func (r *Repository) GetLFUFeeds(ctx context.Context, params *GetLFUFeedsParams) ([]*models.Feed, error) {
	q := `
with current_subscriptions as (
	select distinct feed_id
	from subscriptions
)
select
	feed_id,
	url,
	etag,
	last_modified,
	last_fetched_at,
	created_at
from feeds
join current_subscriptions using (feed_id)
order by last_fetched_at asc
limit $1
 `

	rows, err := r.pool.Query(ctx, q, params.Limit)

	if err != nil {
		return nil, err
	}

	feeds, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.Feed])

	return feeds, err
}
