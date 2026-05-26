package repository

import "context"

func (r *Repository) DeleteSubscription(ctx context.Context, feedID, subscriberID int64) error {
	q := `
delete from subscriptions
where
	feed_id = ? and
	subscriber_id = ?
	`

	_, err := r.db.ExecContext(ctx, q, feedID, subscriberID)

	return err
}
