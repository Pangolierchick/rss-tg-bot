package cron

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type Option func(*Crontab)

type Crontab struct {
	tasks []*Task

	lastRun  time.Time
	lastTick time.Time

	pollingInterval time.Duration
	mu              sync.Mutex
}

func New(opts ...Option) *Crontab {
	tasks := make([]*Task, 0, 16)

	tab := &Crontab{
		tasks:           tasks,
		pollingInterval: 1 * time.Second,
	}

	for _, opt := range opts {
		opt(tab)
	}

	return tab
}

func WithPollingInterval(i time.Duration) Option {
	return func(c *Crontab) {
		c.pollingInterval = i
	}
}

func (c *Crontab) AddTask(expression string, f Func) error {
	task, err := NewTask(expression, f)

	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.tasks = append(c.tasks, task)

	return nil
}

func (c *Crontab) Run(ctx context.Context) chan struct{} {
	ticker := time.NewTicker(c.pollingInterval)
	done := make(chan struct{})

	go func() {
		defer ticker.Stop()
		defer close(done)

		for {
			select {
			case <-ticker.C:
				c.internalRun()
			case <-ctx.Done():
				slog.Warn("Crontab completed")
				return
			}
		}
	}()

	return done
}

func (c *Crontab) internalRun() {
	now := time.Now().Truncate(time.Second)

	c.mu.Lock()
	if c.lastTick.Equal(now) {
		c.mu.Unlock()
		return
	}
	c.lastTick = now
	tasksToRun := make([]*Task, 0)
	for _, t := range c.tasks {
		if t.Match(now) {
			tasksToRun = append(tasksToRun, t)
		}
	}
	c.mu.Unlock()

	for _, t := range tasksToRun {
		go func(t *Task) {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("task panicked:", "reason", r)
				}
			}()
			t.task()
		}(t)
	}
}
