package repository

import (
	"context"
	"fmt"

	"github.com/Pangolierchick/rss-tg-bot/internal/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) AddFeedItems(ctx context.Context, items []*models.FeedItem) error {
	if len(items) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qInsert := `
insert into feed_items (feed_id, guid, title, link, published_at, content_hash)
values ($1, $2, $3, $4, $5, $6)
on conflict do nothing
returning item_id;
	`

	qDispatch := `
insert into deliveries (subscriber_id, feed_item_id)
select
    s.subscriber_id,
    fi.id
from
    subscriptions s
cross join
    unnest($1::bigint[]) as fi(id)
where
    s.feed_id = $2
on conflict do nothing;
	`

	batch := &pgx.Batch{}
	for _, item := range items {
		batch.Queue(qInsert, item.FeedID, item.GUID, item.Title, item.Link, item.PublishedAt, item.ContentHash)
	}

	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	var insertedItemIDs []int64
	for i := 0; i < batch.Len(); i++ {
		row, err := results.Query()
		if err != nil {
			return fmt.Errorf("error in batch statement %d: %w", i, err)
		}

		for row.Next() {
			var itemID int64
			err := row.Scan(&itemID)
			if err != nil {
				row.Close()
				return fmt.Errorf("error scanning item_id: %w", err)
			}
			insertedItemIDs = append(insertedItemIDs, itemID)
		}
		row.Close()
	}

	if err := results.Close(); err != nil {
		return err
	}

	if len(insertedItemIDs) > 0 {
		_, err := tx.Exec(ctx, qDispatch, insertedItemIDs, items[0].FeedID)
		if err != nil {
			return fmt.Errorf("error dispatching deliveries: %w", err)
		}
	}

	return tx.Commit(ctx)
}
