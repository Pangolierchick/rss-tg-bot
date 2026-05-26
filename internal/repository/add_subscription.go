package repository

import (
	"context"

	"github.com/Pangolierchick/rss-tg-bot/internal/models"
)

func (r *Repository) AddSubscription(ctx context.Context, subscription *models.Subscription) error {
	q := `
insert or ignore into subscriptions (feed_id, subscriber_id)
values (?, ?)
	`

	_, err := r.db.ExecContext(ctx, q, subscription.FeedID, subscription.SubscriberID)

	return err
}
