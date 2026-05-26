package repository

import (
	"context"

	"github.com/Pangolierchick/rss-tg-bot/internal/models"
)

func (r *Repository) GetSubscriberByTg(ctx context.Context, tgChatID int64) (*models.Subscriber, error) {
	q := `
select
	subscriber_id,
	tg_chat_id,
	created_at
from subscribers
where
	tg_chat_id = ?
	`

	return scanSubscriber(r.db.QueryRowContext(ctx, q, tgChatID))
}
