package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Pangolierchick/rss-tg-bot/internal/repository"
	"github.com/Pangolierchick/rss-tg-bot/internal/services/fetcher"
	"github.com/Pangolierchick/rss-tg-bot/internal/services/sender"
	"github.com/Pangolierchick/rss-tg-bot/internal/services/subscriptioner"
	telegramfrontend "github.com/Pangolierchick/rss-tg-bot/internal/telegram/frontend"
	tgsender "github.com/Pangolierchick/rss-tg-bot/internal/telegram/sender"
	"github.com/Pangolierchick/rss-tg-bot/pkg/cron"
	"github.com/enetx/surf"
	"github.com/go-telegram/bot"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mmcdole/gofeed"
)

const (
	telegramTokenEnv = "TELEGRAM_TOKEN"
	postgresqlDSN    = "POSTGRESQL_DSN"
	fetchingCron     = "FETCH_CRON"
	sendCron         = "SEND_CRON"
)

type Config struct {
	TelegramToken string
	PostgresqlDSN string
	FetchCron     string
	SendCron      string
}

func readConfig() Config {
	token, ok := os.LookupEnv(telegramTokenEnv)

	if !ok {
		slog.Error("TELEGRAM_TOKEN must be provided")
		os.Exit(1)
	}

	dsn, ok := os.LookupEnv(postgresqlDSN)

	if !ok {
		slog.Error("POSTGRESQL_DSN must be provided")
		os.Exit(1)
	}

	fetchCron, ok := os.LookupEnv(fetchingCron)

	if !ok {
		slog.Error("FETCH_CRON must be provided")
		os.Exit(1)
	}
	sendCron, ok := os.LookupEnv(sendCron)

	if !ok {
		slog.Error("SEND_CRON must be provided")
		os.Exit(1)
	}

	return Config{
		TelegramToken: token,
		PostgresqlDSN: dsn,
		FetchCron:     fetchCron,
		SendCron:      sendCron,
	}
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	config := readConfig()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	poolConfig, err := pgxpool.ParseConfig(config.PostgresqlDSN)
	if err != nil {
		slog.Error("failed to parse config", "error", err)
		return
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		slog.Error("failed to get db pool", "error", err)
		return
	}

	repo := repository.New(pool)

	surfClient := surf.NewClient().
		Builder().
		Impersonate().Chrome().
		Session().
		Build()

	stdClient := surfClient.Std()
	rss := gofeed.NewParser()
	rss.Client = stdClient

	fetchService := fetcher.New(rss, repo, &fetcher.FetcherOpts{
		Limit: 5,
	})

	telegram, err := bot.New(config.TelegramToken)

	if err != nil {
		slog.Error("failed to init telegram bot",
			"error", err,
		)
		return
	}

	tgSender := tgsender.New(telegram)
	senderService := sender.New(repo, tgSender)

	subscriptionerService := subscriptioner.New(repo)
	frontend := telegramfrontend.New(telegram, subscriptionerService)

	crontab := cron.New()
	crontab.AddTask(config.FetchCron, func() {
		slog.Debug("fetching new items")
		err := fetchService.Fetch(ctx)

		if err != nil {
			slog.Error("failed to fetch", "error", err)
		}
	})

	crontab.AddTask(config.SendCron, func() {
		slog.Debug("sending new deliveries")
		err := senderService.SendBatch(ctx, 50)

		if err != nil {
			slog.Error("failed to send batch", "error", err)
		}
	})

	cronWait := crontab.Run(ctx)
	frontend.Run(ctx)

	slog.Info("App started")

	<-signals
	cancel()
	<-cronWait

	slog.Info("Exitting")
}
