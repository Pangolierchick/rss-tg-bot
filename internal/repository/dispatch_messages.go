package repository

import "context"

func (r *Repository) DispatchMessages(ctx context.Context, subscriberID int64) error {
	q := `
with users_subscriptions as (
	select
		feed_id
	from subscriptions
where
		subscriber_id = ?
), items_to_send as (
	select
		item_id
	from feed_items
	where
		feed_id in (select feed_id from users_subscriptions)
	except
	select
		feed_item_id as item_id
	from deliveries
	where
		subscriber_id = ? and
	status = 'sent'
)
insert or ignore into deliveries (subscriber_id, feed_item_id, status)
select ?, item_id, 'pending'
from items_to_send;
	`

	_, err := r.db.ExecContext(ctx, q, subscriberID, subscriberID, subscriberID)

	return err
}
