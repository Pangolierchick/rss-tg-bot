package repository

import (
	"context"

	"github.com/Pangolierchick/rss-tg-bot/internal/models"
)

func (r *Repository) AddSubscription(ctx context.Context, subscription *models.Subscription) error {
	q := `
insert into subscriptions (feed_id, subscriber_id)
values ($1, $2)
on conflict do nothing
	`

	_, err := r.pool.Exec(ctx, q, subscription.FeedID, subscription.SubscriberID)

	return err
}
