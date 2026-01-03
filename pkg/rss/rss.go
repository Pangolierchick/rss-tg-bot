package rss

import (
	"time"
)

var rssDateFormats = []string{
	time.RFC1123Z,
	time.RFC1123,
	"Mon, 02 Jan 2006 15:04 MST",
	"Mon, 02 Jan 06 15:04:05 -0700",
	"Mon, _2 Jan 2006 15:04:05 -0700",
	"Mon, 02 Jan 2006 15:04:05 -0700",
	"2 Jan 2006 15:04:05 -0700",
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02T15:04:05Z07",
	"2006-01-02T15:04:05-07:00",
	"2006-01-02 15:04:05 -0700",
	"2006-01-02",
}

func ParseDate(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	for _, format := range rssDateFormats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, nil
}
