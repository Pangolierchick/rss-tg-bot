package repository

import "context"

func (r *Repository) DeleteSubscription(ctx context.Context, feedID, subscriberID int64) error {
	q := `
delete from subscriptions
where
	feed_id = $1 and
	subscriber_id = $2
	`

	_, err := r.pool.Exec(ctx, q, feedID, subscriberID)

	return err
}
