package repository

import "context"

func (r *Repository) DispatchMessages(ctx context.Context, subscriberID, feedID, limit int64) error {
	if limit <= 0 {
		return nil
	}

	q := `
with items_to_send as (
	select
		item_id
	from feed_items
	where
		feed_id = ?
	order by coalesce(published_at, created_at) desc, item_id desc
	limit ?
)
insert or ignore into deliveries (subscriber_id, feed_item_id, status)
select ?, item_id, 'pending'
from items_to_send;
	`

	_, err := r.db.ExecContext(ctx, q, feedID, limit, subscriberID)

	return err
}
