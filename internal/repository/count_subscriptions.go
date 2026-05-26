package repository

import "context"

func (r *Repository) CountSubscriptions(ctx context.Context, tgChatId int64) (int64, error) {
	q := `
select
	count(*)
from subscribers
join subscriptions using (subscriber_id)
where
	tg_chat_id = ?
	`

	var count int64
	err := r.db.QueryRowContext(ctx, q, tgChatId).Scan(&count)

	return count, err
}
