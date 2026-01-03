package repository

import "context"

func (r *Repository) GetUserSubscriptionsURLs(ctx context.Context, tgChatID int64) ([]string, error) {
	q := `
select
	url
from subscriptions
join subscribers using (subscriber_id)
join feeds using (feed_id)
where
	tg_chat_id = $1
	`

	rows, err := r.pool.Query(ctx, q, tgChatID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	URLs := make([]string, 0, 8)
	for rows.Next() {
		var URL string
		err := rows.Scan(&URL)

		if err != nil {
			return nil, err
		}

		URLs = append(URLs, URL)
	}

	return URLs, rows.Err()
}
