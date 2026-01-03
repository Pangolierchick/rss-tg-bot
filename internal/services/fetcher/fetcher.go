package fetcher

import (
	"context"
	"crypto/sha1"
	"log/slog"

	"github.com/Pangolierchick/rss-tg-bot/internal/models"
	r "github.com/Pangolierchick/rss-tg-bot/internal/repository"
	v2 "github.com/Pangolierchick/rss-tg-bot/pkg/rss/v2"
)

type rssFetcher interface {
	Fetch(URL string) (*v2.Feed, error)
}

type repository interface {
	GetLFUFeeds(ctx context.Context, params *r.GetLFUFeedsParams) ([]*models.Feed, error)
	AddFeedItems(ctx context.Context, items []*models.FeedItem) error
}

type FetcherOpts struct {
	Limit int
}

type Fetcher struct {
	rss  rssFetcher
	repo repository

	opts *FetcherOpts
}

func New(rss rssFetcher, repo repository, opts *FetcherOpts) *Fetcher {
	return &Fetcher{
		rss:  rss,
		repo: repo,
		opts: opts,
	}
}

func (f *Fetcher) Fetch(ctx context.Context) error {
	feeds, err := f.repo.GetLFUFeeds(ctx, &r.GetLFUFeedsParams{Limit: f.opts.Limit})

	if err != nil {
		slog.Error("failed to get least frequently updated feeds", "error", err)
		return err
	}

	for _, feed := range feeds {
		rss, err := f.rss.Fetch(feed.URL)

		if err != nil {
			slog.Error("failed to fetch rss",
				"url", feed.URL,
				"error", err,
			)
			continue
		}

		if len(rss.Channel.Items) == 0 {
			slog.Warn("no new items in channel",
				"url", feed.URL,
			)
			continue
		}

		slog.Debug("Recieved new items. Saving.",
			"url", feed.URL,
			"items", len(rss.Channel.Items),
		)

		modelItems := make([]*models.FeedItem, 0, len(rss.Channel.Items))

		for _, rssItem := range rss.Channel.Items {
			h := sha1.New()
			h.Write([]byte(rssItem.Title + rssItem.Description))
			modelItems = append(modelItems, &models.FeedItem{
				FeedID:      feed.ID,
				GUID:        rssItem.ID,
				Title:       rssItem.Title,
				Link:        rssItem.Link,
				PublishedAt: &rssItem.ParsedPubDate,
				ContentHash: h.Sum(nil),
			})
		}
		err = f.repo.AddFeedItems(ctx, modelItems)

		if err != nil {
			slog.Error("failed to insert new feed items",
				"url", feed.URL,
				"feed_id", feed.ID,
				"error", err,
			)
		}
	}

	return nil
}
