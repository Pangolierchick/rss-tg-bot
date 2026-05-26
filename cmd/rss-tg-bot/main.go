package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Pangolierchick/rss-tg-bot/internal/repository"
	"github.com/Pangolierchick/rss-tg-bot/internal/services/fetcher"
	"github.com/Pangolierchick/rss-tg-bot/internal/services/sender"
	"github.com/Pangolierchick/rss-tg-bot/internal/services/subscriptioner"
	telegramfrontend "github.com/Pangolierchick/rss-tg-bot/internal/telegram/frontend"
	tgsender "github.com/Pangolierchick/rss-tg-bot/internal/telegram/sender"
	"github.com/go-telegram/bot"
	"github.com/mmcdole/gofeed"
	"gopkg.in/yaml.v3"
	_ "modernc.org/sqlite"
)

const (
	defaultConfigPath = "config.yaml"
)

type Config struct {
	Telegram TelegramConfig `yaml:"telegram"`
	Database DatabaseConfig `yaml:"database"`
	Fetch    FetchConfig    `yaml:"fetch"`
	Send     SendConfig     `yaml:"send"`
}

type TelegramConfig struct {
	Token    string `yaml:"token"`
	ProxyURL string `yaml:"proxy_url"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type FetchConfig struct {
	Interval string        `yaml:"interval"`
	Limit    int           `yaml:"limit"`
	ProxyURL string        `yaml:"proxy_url"`
	Every    time.Duration `yaml:"-"`
}

type SendConfig struct {
	Interval  string        `yaml:"interval"`
	BatchSize int64         `yaml:"batch_size"`
	Every     time.Duration `yaml:"-"`
}

func readConfig(path string) (Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(raw, &config); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}

	if config.Telegram.Token == "" {
		return Config{}, fmt.Errorf("telegram.token must be provided")
	}
	if config.Database.Path == "" {
		return Config{}, fmt.Errorf("database.path must be provided")
	}
	if config.Fetch.Interval == "" {
		return Config{}, fmt.Errorf("fetch.interval must be provided")
	}
	if config.Send.Interval == "" {
		return Config{}, fmt.Errorf("send.interval must be provided")
	}

	config.Fetch.Every, err = time.ParseDuration(config.Fetch.Interval)
	if err != nil {
		return Config{}, fmt.Errorf("parse fetch.interval: %w", err)
	}
	config.Send.Every, err = time.ParseDuration(config.Send.Interval)
	if err != nil {
		return Config{}, fmt.Errorf("parse send.interval: %w", err)
	}
	if config.Fetch.Every <= 0 {
		return Config{}, fmt.Errorf("fetch.interval must be positive")
	}
	if config.Send.Every <= 0 {
		return Config{}, fmt.Errorf("send.interval must be positive")
	}
	if config.Fetch.Limit <= 0 {
		config.Fetch.Limit = 5
	}
	if config.Send.BatchSize <= 0 {
		config.Send.BatchSize = 50
	}

	return config, nil
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	configPath := flag.String("config", defaultConfigPath, "path to YAML config file")
	flag.Parse()

	config, err := readConfig(*configPath)
	if err != nil {
		slog.Error("failed to read config", "error", err)
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	db, err := sql.Open("sqlite", sqliteDSN(config.Database.Path))
	if err != nil {
		slog.Error("failed to open database", "error", err)
		return
	}
	defer db.Close()
	db.SetMaxOpenConns(1)

	repo := repository.New(db)
	if err := repo.Init(ctx); err != nil {
		slog.Error("failed to initialize database", "error", err)
		return
	}

	stdClient, err := httpClient(30*time.Second, config.Fetch.ProxyURL)
	if err != nil {
		slog.Error("failed to configure fetch http client",
			"error", err,
		)
		return
	}
	rss := gofeed.NewParser()

	fetchService := fetcher.New(rss, stdClient, repo, &fetcher.FetcherOpts{
		Limit: config.Fetch.Limit,
	})

	telegramOpts := make([]bot.Option, 0, 1)
	if config.Telegram.ProxyURL != "" {
		telegramHTTPClient, err := httpClient(time.Minute, config.Telegram.ProxyURL)
		if err != nil {
			slog.Error("failed to configure telegram proxy",
				"error", err,
			)
			return
		}
		telegramOpts = append(telegramOpts, bot.WithHTTPClient(time.Minute, telegramHTTPClient))
	}

	telegram, err := bot.New(config.Telegram.Token, telegramOpts...)

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

	fetchDone := runEvery(ctx, config.Fetch.Every, func(ctx context.Context) {
		slog.Debug("fetching new items")
		err := fetchService.Fetch(ctx)

		if err != nil {
			slog.Error("failed to fetch", "error", err)
		}
	})

	sendDone := runEvery(ctx, config.Send.Every, func(ctx context.Context) {
		slog.Debug("sending new deliveries")
		err := senderService.SendBatch(ctx, config.Send.BatchSize)

		if err != nil {
			slog.Error("failed to send batch", "error", err)
		}
	})

	frontend.Run(ctx)

	slog.Info("App started")

	<-signals
	cancel()
	<-fetchDone
	<-sendDone

	slog.Info("Exitting")
}

func sqliteDSN(path string) string {
	return path + "?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)"
}

func httpClient(timeout time.Duration, proxyURL string) (*http.Client, error) {
	client := &http.Client{Timeout: timeout}
	if proxyURL == "" {
		return client, nil
	}

	parsedProxyURL, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("parse proxy url: %w", err)
	}

	client.Transport = &http.Transport{
		Proxy: http.ProxyURL(parsedProxyURL),
	}

	return client, nil
}

func runEvery(ctx context.Context, interval time.Duration, task func(context.Context)) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)

		timer := time.NewTimer(0)
		defer timer.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				func() {
					defer func() {
						if r := recover(); r != nil {
							slog.Error("timer task panicked", "reason", r)
						}
					}()

					task(ctx)
				}()
				timer.Reset(interval)
			}
		}
	}()

	return done
}
