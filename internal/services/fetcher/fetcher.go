package fetcher

import (
	"context"
	"crypto/sha1"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Pangolierchick/rss-tg-bot/internal/models"
	r "github.com/Pangolierchick/rss-tg-bot/internal/repository"
	"github.com/mmcdole/gofeed"
)

type repository interface {
	GetLFUFeeds(ctx context.Context, params *r.GetLFUFeedsParams) ([]*models.Feed, error)
	AddFeedItems(ctx context.Context, items []*models.FeedItem) error
	UpdateFeed(ctx context.Context, feed *models.Feed) error
}

type FetcherOpts struct {
	Limit int
}

type Fetcher struct {
	rss    *gofeed.Parser
	client *http.Client
	repo   repository

	opts *FetcherOpts
}

func New(rss *gofeed.Parser, client *http.Client, repo repository, opts *FetcherOpts) *Fetcher {
	return &Fetcher{
		rss:    rss,
		client: client,
		repo:   repo,
		opts:   opts,
	}
}

func (f *Fetcher) Fetch(ctx context.Context) error {
	feeds, err := f.repo.GetLFUFeeds(ctx, &r.GetLFUFeedsParams{Limit: f.opts.Limit})

	if err != nil {
		slog.Error("failed to get least frequently updated feeds", "error", err)
		return err
	}

	for _, feed := range feeds {
		rss, err := f.fetchFeed(ctx, feed)

		if err != nil {
			slog.Error("failed to fetch rss",
				"url", feed.URL,
				"error", err,
			)
			continue
		}
		if rss == nil {
			continue
		}

		if len(rss.Items) == 0 {
			slog.Warn("no new items in channel",
				"url", feed.URL,
			)
			continue
		}

		slog.Debug("Recieved new items. Saving.",
			"url", feed.URL,
			"items", len(rss.Items),
		)

		modelItems := make([]*models.FeedItem, 0, len(rss.Items))

		for _, rssItem := range rss.Items {
			h := sha1.New()
			h.Write([]byte(rssItem.Title + rssItem.Description))
			hash := h.Sum(nil)

			guid := rssItem.GUID

			if len(guid) == 0 {
				guid = rssItem.Link
			}

			if len(guid) == 0 {
				guid = string(hash)
			}

			modelItems = append(modelItems, &models.FeedItem{
				FeedID:      feed.ID,
				GUID:        guid,
				Title:       rssItem.Title,
				Link:        rssItem.Link,
				PublishedAt: rssItem.PublishedParsed,
				ContentHash: hash,
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

func (f *Fetcher) fetchFeed(ctx context.Context, feed *models.Feed) (*gofeed.Feed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feed.URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/rss+xml, application/atom+xml, application/xml;q=0.9, text/xml;q=0.8, */*;q=0.1")
	req.Header.Set("User-Agent", "rss-tg-bot/1.0 (+https://github.com/Pangolierchick/rss-tg-bot)")

	if feed.ETag != nil && *feed.ETag != "" {
		req.Header.Set("If-None-Match", *feed.ETag)
	}
	if feed.LastModified != nil {
		req.Header.Set("If-Modified-Since", feed.LastModified.UTC().Format(http.TimeFormat))
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	now := time.Now().UTC()
	if resp.StatusCode == http.StatusNotModified {
		feed.LastFetchedAt = now
		return nil, f.repo.UpdateFeed(ctx, feed)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	parsed, err := f.rss.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	if etag := resp.Header.Get("ETag"); etag != "" {
		feed.ETag = &etag
	}
	if lastModified := resp.Header.Get("Last-Modified"); lastModified != "" {
		if t, err := http.ParseTime(lastModified); err == nil {
			feed.LastModified = &t
		}
	}
	feed.LastFetchedAt = now

	if err := f.repo.UpdateFeed(ctx, feed); err != nil {
		return nil, err
	}

	return parsed, nil
}
