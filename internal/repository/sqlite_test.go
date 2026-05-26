package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Pangolierchick/rss-tg-bot/internal/models"
	_ "modernc.org/sqlite"
)

func TestDispatchMessagesSQLite(t *testing.T) {
	ctx := context.Background()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)

	repo := New(db)
	if err := repo.Init(ctx); err != nil {
		t.Fatal(err)
	}

	feedID, err := repo.AddFeed(ctx, &models.Feed{URL: "https://example.com/rss.xml"})
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.AddFeedItems(ctx, []*models.FeedItem{
		{
			FeedID:      feedID,
			GUID:        "guid-1",
			Title:       "Title",
			Link:        "https://example.com/post",
			ContentHash: []byte("hash-1"),
		},
	}); err != nil {
		t.Fatal(err)
	}

	subscriberID, err := repo.AddSubscriber(ctx, &models.Subscriber{TgChatID: 123})
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.AddSubscription(ctx, &models.Subscription{FeedID: feedID, SubscriberID: subscriberID}); err != nil {
		t.Fatal(err)
	}
	if err := repo.DispatchMessages(ctx, subscriberID); err != nil {
		t.Fatal(err)
	}

	tx, err := repo.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	deliveries, err := repo.GetDeliveries(ctx, tx, &GetDeliveriesParams{
		Status: models.DeliveryStatusPending,
		Limit:  10,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(deliveries) != 1 {
		t.Fatalf("expected 1 delivery, got %d", len(deliveries))
	}
	if deliveries[0].SubscriberID != subscriberID {
		t.Fatalf("unexpected subscriber id: %d", deliveries[0].SubscriberID)
	}
}
