package repository

import (
	"context"
	"database/sql"

	"github.com/Pangolierchick/rss-tg-bot/internal/models"
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
limit ?
 `

	rows, err := r.db.QueryContext(ctx, q, params.Limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feeds := make([]*models.Feed, 0, params.Limit)
	for rows.Next() {
		feed, err := scanFeedRows(rows)
		if err != nil {
			return nil, err
		}

		feeds = append(feeds, feed)
	}

	return feeds, rows.Err()
}

func scanFeedRows(rows *sql.Rows) (*models.Feed, error) {
	var feed models.Feed
	var lastModified sql.NullInt64
	var lastFetchedAt int64
	var createdAt int64

	if err := rows.Scan(
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
