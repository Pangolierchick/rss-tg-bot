package repository

import (
	"context"
	"database/sql"

	"github.com/Pangolierchick/rss-tg-bot/internal/models"
)

func (r *Repository) GetSubscriber(ctx context.Context, subscriberID int64) (*models.Subscriber, error) {
	q := `
select
	subscriber_id,
	tg_chat_id,
	created_at
from subscribers
where
	subscriber_id = ?
	`

	return scanSubscriber(r.db.QueryRowContext(ctx, q, subscriberID))
}

func scanSubscriber(row *sql.Row) (*models.Subscriber, error) {
	var subscriber models.Subscriber
	var createdAt int64

	if err := row.Scan(
		&subscriber.SubscriberID,
		&subscriber.TgChatID,
		&createdAt,
	); err != nil {
		return nil, err
	}

	subscriber.CreatedAt = scanUnixTime(createdAt)

	return &subscriber, nil
}
