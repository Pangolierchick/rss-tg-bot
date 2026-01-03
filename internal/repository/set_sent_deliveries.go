package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) SetSentStatusDeliveries(ctx context.Context, tx pgx.Tx, deliveryIDs []int64) error {
	q := `
update deliveries
set
	status = 'sent',
	sent_at = now()
where
	delivery_id = any($1)
	`

	_, err := tx.Exec(ctx, q, deliveryIDs)

	return err
}
