package repository

import (
	"context"

	"github.com/Pangolierchick/rss-tg-bot/internal/models"
	"github.com/jackc/pgx/v5"
)

type GetDeliveriesParams struct {
	Status models.DeliveryStatus
	Limit  int64
}

func (r *Repository) GetDeliveries(ctx context.Context, tx pgx.Tx, params *GetDeliveriesParams) ([]*models.Delivery, error) {
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
    d.status = $1
order by f.published_at
limit $2;
	`
	rows, err := tx.Query(ctx, q, params.Status, params.Limit)

	if err != nil {
		return nil, err
	}

	deliveries, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.Delivery])

	return deliveries, err
}
