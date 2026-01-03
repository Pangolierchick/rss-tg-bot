package repository

import (
	"context"

	"github.com/Pangolierchick/rss-tg-bot/internal/models"
)

func (r *Repository) AddFeed(ctx context.Context, feed *models.Feed) (int64, error) {
	q := `
insert into feeds (url)
values ($1)
on conflict (url) do update
set
	last_modified = now()
returning feed_id
	`

	var feedID int64
	err := r.pool.QueryRow(ctx, q, feed.URL).Scan(&feedID)

	return feedID, err
}
