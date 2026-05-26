package repository

import (
	"context"
	"database/sql"
	"time"
)

func (r *Repository) SetSentStatusDeliveries(ctx context.Context, tx *sql.Tx, deliveryIDs []int64) error {
	if len(deliveryIDs) == 0 {
		return nil
	}

	q := `
update deliveries
set
	status = 'sent',
	sent_at = ?
where
	delivery_id in (` + placeholders(len(deliveryIDs)) + `)
	`

	args := make([]any, 0, len(deliveryIDs)+1)
	args = append(args, unixTime(time.Now()))
	for _, id := range deliveryIDs {
		args = append(args, id)
	}

	_, err := tx.ExecContext(ctx, q, args...)

	return err
}
