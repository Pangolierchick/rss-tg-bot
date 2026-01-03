package model

import "time"

type Subscription struct {
	ID         int64      `db:"subscription_id"`
	UserID     int64      `db:"user_id"`
	URL        string     `db:"url"`
	LastPolled *time.Time `db:"last_polled"`
	CreatedAt  time.Time  `db:"created_at"`
}
