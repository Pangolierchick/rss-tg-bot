package repository

import (
	"context"

	"github.com/Pangolierchick/rss-tg-bot/internal/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetFeedByURL(ctx context.Context, URL string) (*models.Feed, error) {
	q := `
select
	feed_id,
	url,
	etag,
	last_modified,
	last_fetched_at,
	created_at
from feeds
where
	url = $1
	`

	rows, err := r.pool.Query(ctx, q, URL)

	if err != nil {
		return nil, err
	}

	feed, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[models.Feed])

	return feed, err
}
