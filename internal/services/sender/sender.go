package sender

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Pangolierchick/rss-tg-bot/internal/models"
	r "github.com/Pangolierchick/rss-tg-bot/internal/repository"
	"github.com/jackc/pgx/v5"
)

type messageSender interface {
	Send(ctx context.Context, ID any, message string) error
}

type repository interface {
	SetSentStatusDeliveries(ctx context.Context, tx pgx.Tx, deliveryIDs []int64) error
	GetDeliveries(ctx context.Context, tx pgx.Tx, params *r.GetDeliveriesParams) ([]*models.Delivery, error)
	Begin(ctx context.Context) (pgx.Tx, error)
	GetSubscriber(ctx context.Context, subscriberID int64) (*models.Subscriber, error)
	GetFeedItem(ctx context.Context, feedItemID int64) (*models.FeedItem, error)
}

type Service struct {
	repo   repository
	sender messageSender
}

func New(repo repository, sender messageSender) *Service {
	return &Service{
		repo:   repo,
		sender: sender,
	}
}

func (s *Service) SendBatch(ctx context.Context, limit int64) error {
	tx, err := s.repo.Begin(ctx)

	if err != nil {
		slog.Error("failed to begin transaction",
			"error", err,
		)
		return err
	}

	defer tx.Rollback(ctx)

	deliveries, err := s.repo.GetDeliveries(ctx, tx, &r.GetDeliveriesParams{
		Status: models.DeliveryStatusPending,
		Limit:  limit,
	})

	if err != nil {
		slog.Error("failed to get pending deliveries",
			"error", err,
		)
		return err
	}

	successedDeliveryIDs := make([]int64, 0, len(deliveries))

	for _, delivery := range deliveries {
		sub, err := s.repo.GetSubscriber(ctx, delivery.SubscriberID)

		if err != nil {
			slog.Error("failed to get subscriber to deliver message",
				"subscriber_id", delivery.SubscriberID,
				"error", err,
			)
			continue
		}

		item, err := s.repo.GetFeedItem(ctx, delivery.FeedItemID)

		if err != nil {
			slog.Error("failed to get feed item",
				"item_id", delivery.FeedItemID,
				"subscriber_id", delivery.SubscriberID,
				"error", err,
			)
			continue
		}

		msg := collectMessage(item)
		err = s.sender.Send(ctx, sub.TgChatID, msg)

		if err != nil {
			slog.Error("failed to deliver message",
				"chat_id", sub.TgChatID,
				"error", err,
				"message", msg,
			)
			continue
		}

		successedDeliveryIDs = append(successedDeliveryIDs, delivery.ID)
	}

	err = s.repo.SetSentStatusDeliveries(ctx, tx, successedDeliveryIDs)

	if err != nil {
		slog.Error("failed to set sent status on deliveries",
			"error", err,
		)

		return err
	}

	return tx.Commit(ctx)
}

func collectMessage(f *models.FeedItem) string {
	return fmt.Sprintf("*%s*\n\n%s", f.Title, f.Link)
}
