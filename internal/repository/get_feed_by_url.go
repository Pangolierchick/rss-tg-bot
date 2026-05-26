package repository

import (
	"context"
	"database/sql"

	"github.com/Pangolierchick/rss-tg-bot/internal/models"
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
	url = ?
	`

	return scanFeed(r.db.QueryRowContext(ctx, q, URL))
}

func scanFeed(row *sql.Row) (*models.Feed, error) {
	var feed models.Feed
	var lastModified sql.NullInt64
	var lastFetchedAt int64
	var createdAt int64

	if err := row.Scan(
		&feed.ID,
		&feed.URL,
		&feed.ETag,
		&lastModified,
		&lastFetchedAt,
		&createdAt,
	); err != nil {
		return nil, err
	}

	feed.LastModified = scanNullUnixTime(lastModified)
	feed.LastFetchedAt = scanUnixTime(lastFetchedAt)
	feed.CreatedAt = scanUnixTime(createdAt)

	return &feed, nil
}
