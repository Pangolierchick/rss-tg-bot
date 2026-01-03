package repository

import "context"

func (r *Repository) CountSubscriptions(ctx context.Context, tgChatId int64) (int64, error) {
	q := `
select
	count(*)
from subscribers
join subscriptions using (subscriber_id)
where
	tg_chat_id = $1
	`

	var count int64
	err := r.pool.QueryRow(ctx, q, tgChatId).Scan(&count)

	return count, err
}
