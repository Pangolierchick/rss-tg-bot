package repository

import (
	"context"
	"database/sql"

	"github.com/Pangolierchick/rss-tg-bot/internal/models"
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
	item_id = ?
	`

	return scanFeedItem(r.db.QueryRowContext(ctx, q, feedItemID))
}

func scanFeedItem(row *sql.Row) (*models.FeedItem, error) {
	var item models.FeedItem
	var publishedAt sql.NullInt64
	var createdAt int64

	if err := row.Scan(
		&item.ID,
		&item.FeedID,
		&item.GUID,
		&item.Title,
		&item.Link,
		&publishedAt,
		&item.ContentHash,
		&createdAt,
	); err != nil {
		return nil, err
	}

	item.PublishedAt = scanNullUnixTime(publishedAt)
	item.CreatedAt = scanUnixTime(createdAt)

	return &item, nil
}
