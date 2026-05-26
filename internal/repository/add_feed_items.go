package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Pangolierchick/rss-tg-bot/internal/models"
)

func (r *Repository) AddFeedItems(ctx context.Context, items []*models.FeedItem) error {
	if len(items) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	qInsert := `
insert or ignore into feed_items (feed_id, guid, title, link, published_at, content_hash)
values (?, ?, ?, ?, ?, ?)
returning item_id;
	`

	qDispatch := `
insert or ignore into deliveries (subscriber_id, feed_item_id)
select
    s.subscriber_id,
    ?
from subscriptions s
where
    s.feed_id = ?;
	`

	for _, item := range items {
		var itemID int64
		err := tx.QueryRowContext(
			ctx,
			qInsert,
			item.FeedID,
			item.GUID,
			item.Title,
			item.Link,
			unixTimePtr(item.PublishedAt),
			item.ContentHash,
		).Scan(&itemID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}

			return fmt.Errorf("insert feed item: %w", err)
		}

		if _, err := tx.ExecContext(ctx, qDispatch, itemID, item.FeedID); err != nil {
			return fmt.Errorf("dispatch deliveries: %w", err)
		}
	}

	return tx.Commit()
}
