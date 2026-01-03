package repository

import (
	"context"

	"github.com/Pangolierchick/rss-tg-bot/internal/models"
)

func (r *Repository) UpdateFeed(ctx context.Context, feed *models.Feed) error {
	q := `
update feeds
set
	url = $2,
	etag = $3,
	last_modified = $4,
	last_fetched_at = $5,
	created_at = $6
where
	feed_id = $1
	`

	_, err := r.pool.Exec(ctx, q,
		feed.ID,
		feed.URL,
		feed.ETag,
		feed.LastModified,
		feed.LastFetchedAt,
		feed.CreatedAt,
	)

	return err
}
