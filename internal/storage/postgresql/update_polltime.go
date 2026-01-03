package postgresql

import "context"

func (r *Repository) UpdatePolltime(ctx context.Context, IDs []int64) error {
	q := `
update subscriptions
set
	last_polled = now()
where
	subscription_id = any($1)
	`

	_, err := r.pool.Exec(ctx, q, IDs)

	return err
}
