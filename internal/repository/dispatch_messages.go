package repository

import "context"

func (r *Repository) DispatchMessages(ctx context.Context, subscriberID int64) error {
	q := `
with users_subscriptions as (
	select
		feed_id
	from subscriptions
	where
		subscriber_id = $1
), items_to_send as (
	select
		item_id
	from feed_items
	where
		feed_id = any((select * from users_subscriptions))
	except
	select
		feed_item_id as item_id
	from deliveries
	where
		subscriber_id = $1 and
	status = 'sent'
)
insert into deliveries (subscriber_id, feed_item_id, status)
select $1, item_id, 'pending'
from items_to_send
on conflict do nothing;
	`

	_, err := r.pool.Exec(ctx, q, subscriberID)

	return err
}
