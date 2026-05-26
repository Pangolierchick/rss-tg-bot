package repository

import (
	"context"
	"database/sql"

	"github.com/Pangolierchick/rss-tg-bot/internal/models"
)

type GetDeliveriesParams struct {
	Status models.DeliveryStatus
	Limit  int64
}

func (r *Repository) GetDeliveries(ctx context.Context, tx *sql.Tx, params *GetDeliveriesParams) ([]*models.Delivery, error) {
	q := `
select
    d.delivery_id,
    d.subscriber_id,
    d.feed_item_id,
    d.status,
    d.sent_at,
    d.created_at
from deliveries d
join feed_items f on f.item_id = d.feed_item_id
where
    d.status = ?
order by coalesce(f.published_at, f.created_at), d.created_at
limit ?;
	`
	rows, err := tx.QueryContext(ctx, q, params.Status, params.Limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deliveries := make([]*models.Delivery, 0, params.Limit)
	for rows.Next() {
		var delivery models.Delivery
		var sentAt sql.NullInt64
		var createdAt int64

		if err := rows.Scan(
			&delivery.ID,
			&delivery.SubscriberID,
			&delivery.FeedItemID,
			&delivery.Status,
			&sentAt,
			&createdAt,
		); err != nil {
			return nil, err
		}

		delivery.SentAt = scanNullUnixTime(sentAt)
		delivery.CreatedAt = scanUnixTime(createdAt)
		deliveries = append(deliveries, &delivery)
	}

	return deliveries, rows.Err()
}
