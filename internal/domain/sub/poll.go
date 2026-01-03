package sub

import (
	"context"
	"log/slog"
)

func (s *Subscriptioner) Poll(ctx context.Context) error {
	slog.Debug("polling subscriptions")

	subs, err := s.repo.GetSubscriptions(ctx)

	if err != nil {
		slog.Error("failed to get subscriptions",
			"error", err,
		)
		return err
	}

	successed := make([]int64, 0, len(subs))

	for _, sub := range subs {
		feed, err := s.rss.Fetch(sub.URL)

		if err != nil {
			slog.Error("failed to fetch rss",
				"url", sub.URL,
				"userID", sub.UserID,
				"error", err,
			)
			continue
		}

		err = s.Send(ctx, sub.UserID, sub.ID, feed.Channel.Items)

		if err != nil {
			slog.Error("failed to send rss items",
				"error", err,
				"userID", sub.UserID,
				"url", sub.URL,
			)
		}

		successed = append(successed, sub.ID)
	}

	err = s.repo.UpdatePolltime(ctx, successed)

	return err
}
