package models

import "time"

type Feed struct {
	ID            int64      `db:"feed_id" json:"feed_id"`
	URL           string     `db:"url" json:"url"`
	ETag          *string    `db:"etag" json:"etag,omitempty"`
	LastModified  *time.Time `db:"last_modified" json:"last_modified,omitempty"`
	LastFetchedAt time.Time  `db:"last_fetched_at" json:"last_fetched_at"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
}
type Subscriber struct {
	SubscriberID int64     `db:"subscriber_id" json:"subscriber_id"`
	TgChatID     int64     `db:"tg_chat_id" json:"tg_chat_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type Subscription struct {
	FeedID    int64     `db:"feed_id" json:"feed_id"`
	ID        int64     `db:"subscriber_id" json:"subscriber_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type FeedItem struct {
	ID          int64      `db:"item_id" json:"item_id"`
	FeedID      int64      `db:"feed_id" json:"feed_id"`
	GUID        string     `db:"guid" json:"guid"`
	Title       string     `db:"title" json:"title"`
	Link        string     `db:"link" json:"link"`
	PublishedAt *time.Time `db:"published_at" json:"published_at,omitempty"`
	ContentHash []byte     `db:"content_hash" json:"content_hash"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
}

type DeliveryStatus string

const (
	DeliveryStatusPending DeliveryStatus = "pending"
	DeliveryStatusSent    DeliveryStatus = "sent"
)

type Delivery struct {
	ID           int64          `db:"delivery_id" json:"delivery_id"`
	SubscriberID int64          `db:"subscriber_id" json:"subscriber_id"`
	FeedItemID   int64          `db:"feed_item_id" json:"feed_item_id"`
	Status       DeliveryStatus `db:"status" json:"status"`
	SentAt       *time.Time     `db:"sent_at" json:"sent_at,omitempty"`
	CreatedAt    time.Time      `db:"created_at" json:"created_at"`
}
