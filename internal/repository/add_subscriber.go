package repository

import (
	"context"

	"github.com/Pangolierchick/rss-tg-bot/internal/models"
)

func (r *Repository) AddSubscriber(ctx context.Context, subscriber *models.Subscriber) (int64, error) {
	q := `
insert into subscribers (tg_chat_id)
values ($1)
on conflict (tg_chat_id) do update
set
	tg_chat_id = excluded.tg_chat_id
returning subscriber_id
	`

	var subscriberID int64
	err := r.pool.QueryRow(ctx, q, subscriber.TgChatID).Scan(&subscriberID)

	return subscriberID, err
}
