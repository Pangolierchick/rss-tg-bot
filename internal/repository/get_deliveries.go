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
	delivery_id,
	subscriber_id,
	feed_item_id,
	status,
	sent_at,
	created_at
from deliveries
where
	status = $1
for update skip locked
limit $2
	`
	rows, err := tx.Query(ctx, q, params.Status, params.Limit)

	if err != nil {
		return nil, err
	}

	deliveries, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.Delivery])

	return deliveries, err
}
