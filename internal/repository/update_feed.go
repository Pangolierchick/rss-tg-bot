package repository

import (
	"context"

	"github.com/Pangolierchick/rss-tg-bot/internal/models"
)

func (r *Repository) UpdateFeed(ctx context.Context, feed *models.Feed) error {
	q := `
update feeds
set
	url = ?,
	etag = ?,
	last_modified = ?,
	last_fetched_at = ?
where
	feed_id = ?
	`

	_, err := r.db.ExecContext(ctx, q,
		feed.URL,
		feed.ETag,
		unixTimePtr(feed.LastModified),
		unixTime(feed.LastFetchedAt),
		feed.ID,
	)

	return err
}
