package model

import "time"

type Item struct {
	ID             string    `db:"item_id"`
	SubscriptionID int64     `db:"subscription_id"`
	CreatedAt      time.Time `db:"created_at"`
}
