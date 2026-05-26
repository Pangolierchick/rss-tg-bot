package sender

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Pangolierchick/rss-tg-bot/internal/models"
	r "github.com/Pangolierchick/rss-tg-bot/internal/repository"
	_ "modernc.org/sqlite"
)

type fakeSender struct {
	sent int
}

func (f *fakeSender) Send(ctx context.Context, ID any, message string) error {
	f.sent++
	return nil
}

func TestSendBatchMarksDeliveriesSentWithSingleSQLiteConnection(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)

	repo := r.New(db)
	if err := repo.Init(ctx); err != nil {
		t.Fatal(err)
	}

	feedID, err := repo.AddFeed(ctx, &models.Feed{URL: "https://example.com/rss.xml"})
	if err != nil {
		t.Fatal(err)
	}
	subscriberID, err := repo.AddSubscriber(ctx, &models.Subscriber{TgChatID: 123})
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.AddSubscription(ctx, &models.Subscription{FeedID: feedID, SubscriberID: subscriberID}); err != nil {
		t.Fatal(err)
	}
	if err := repo.AddFeedItems(ctx, []*models.FeedItem{
		{
			FeedID:      feedID,
			GUID:        "guid-1",
			Title:       "Title * with _ markdown",
			Link:        "https://example.com/post",
			ContentHash: []byte("hash-1"),
		},
	}); err != nil {
		t.Fatal(err)
	}

	sender := &fakeSender{}
	service := New(repo, sender)
	if err := service.SendBatch(ctx, 10); err != nil {
		t.Fatal(err)
	}
	if sender.sent != 1 {
		t.Fatalf("expected 1 sent message, got %d", sender.sent)
	}

	tx, err := repo.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	deliveries, err := repo.GetDeliveries(ctx, tx, &r.GetDeliveriesParams{
		Status: models.DeliveryStatusPending,
		Limit:  10,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(deliveries) != 0 {
		t.Fatalf("expected no pending deliveries, got %d", len(deliveries))
	}
}
