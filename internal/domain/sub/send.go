package sub

import (
	"context"
	"log/slog"

	"github.com/Pangolierchick/rss-tg-bot/internal/model"
	v2 "github.com/Pangolierchick/rss-tg-bot/pkg/rss/v2"
)

func (s *Subscriptioner) Send(ctx context.Context, userID, subscriptionID int64, items []v2.Item) error {
	for _, item := range items {
		ok, err := s.dedup.IsNew(ctx, subscriptionID, item)

		if err != nil {
			slog.Error("failed to dedup.IsNew", "error", err)
			return err
		}

		if ok == true {
			err := s.producer.Send(ctx, userID, item.Title, item.Description, item.ID)
			if err != nil {
				slog.Error("failed to send item",
					"error", err,
				)
			}

			err = s.repo.AddItem(ctx, model.Item{
				ID:             item.ID,
				SubscriptionID: subscriptionID,
			})

			if err != nil {
				slog.Error("failed to add item", "error", err)
			}
		}
	}

	return nil
}
