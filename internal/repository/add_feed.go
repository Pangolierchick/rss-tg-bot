package repository

import (
	"context"

	"github.com/Pangolierchick/rss-tg-bot/internal/models"
)

func (r *Repository) AddFeed(ctx context.Context, feed *models.Feed) (int64, error) {
	q := `
insert into feeds (url)
values (?)
on conflict (url) do update
set
	url = excluded.url
returning feed_id
	`

	var feedID int64
	err := r.db.QueryRowContext(ctx, q, feed.URL).Scan(&feedID)

	return feedID, err
}
