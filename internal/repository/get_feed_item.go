package repository

import (
	"context"

	"github.com/Pangolierchick/rss-tg-bot/internal/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetFeedItem(ctx context.Context, feedItemID int64) (*models.FeedItem, error) {
	q := `
select
	item_id,
	feed_id,
	guid,
	title,
	link,
	published_at,
	content_hash,
	created_at
from feed_items
where
	item_id = $1
	`

	rows, err := r.pool.Query(ctx, q, feedItemID)

	if err != nil {
		return nil, err
	}

	item, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[models.FeedItem])

	return item, err
}
