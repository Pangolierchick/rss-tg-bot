package v2

import (
	"encoding/xml"
	"fmt"
	"io"
	"time"

	"github.com/Pangolierchick/rss-tg-bot/pkg/rss"
	"golang.org/x/net/html/charset"
)

func Parse(r io.Reader) (*Feed, error) {
	var feed Feed
	p := xml.NewDecoder(r)
	p.CharsetReader = charset.NewReaderLabel
	err := p.Decode(&feed)
	if err != nil {
		return nil, err
	}

	if feed.Channel == nil {
		return nil, fmt.Errorf("no channel found in RSS 2.0 feed")
	}

	return &feed, nil
}

func (i *Item) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type Alias Item
	aux := &struct {
		PubDate string `xml:"pubDate"`
		Date    string `xml:"date"`
		*Alias
	}{
		Alias: (*Alias)(i),
	}

	if err := d.DecodeElement(aux, &start); err != nil {
		return err
	}

	if aux.PubDate != "" {
		if t, err := rss.ParseDate(aux.PubDate); err == nil {
			i.ParsedPubDate = t
		} else {
			fmt.Printf("failed to parse pubdate: %v", err)
		}
	}

	if aux.Date != "" {
		if t, err := rss.ParseDate(aux.Date); err == nil {
			i.ParsedDate = t
		} else {
			fmt.Printf("failed to parse pubdate: %v", err)
		}
	}

	return nil
}

type (
	Feed struct {
		XMLName xml.Name `xml:"rss"`
		Channel *Channel `xml:"channel"`
	}

	Channel struct {
		XMLName     xml.Name   `xml:"channel"`
		Title       string     `xml:"title"`
		Language    string     `xml:"language"`
		Author      string     `xml:"author"`
		Description string     `xml:"description"`
		Link        []Link     `xml:"link"`
		Image       *Image     `xml:"image"`
		Categories  []Category `xml:"category"`
		Items       []Item     `xml:"item"`
		MinsToLive  int        `xml:"ttl"`
		SkipHours   []int      `xml:"skipHours>hour"`
		SkipDays    []string   `xml:"skipDays>day"`
	}

	Link struct {
		Rel      string `xml:"rel,attr"`
		Href     string `xml:"href,attr"`
		Type     string `xml:"type,attr"`
		Chardata string `xml:",chardata"`
	}

	Image struct {
		XMLName xml.Name `xml:"image"`
		Href    string   `xml:"href,attr"`
		Title   string   `xml:"title"`
		URL     string   `xml:"url"`
		Height  int      `xml:"height"`
		Width   int      `xml:"width"`
	}

	Category struct {
		XMLName xml.Name `xml:"category"`
		Name    string   `xml:"text,attr"`
	}

	Item struct {
		XMLName       xml.Name    `xml:"item"`
		Title         string      `xml:"title"`
		Description   string      `xml:"description"`
		Content       string      `xml:"encoded"`
		Categories    []string    `xml:"category"`
		Link          string      `xml:"link"`
		PubDate       string      `xml:"pubDate"`
		Date          string      `xml:"date"`
		ParsedPubDate time.Time   `xml:"-"`
		ParsedDate    time.Time   `xml:"-"`
		Image         *Image      `xml:"image"`
		ID            string      `xml:"guid"`
		Enclosures    []Enclosure `xml:"enclosure"`
	}

	Enclosure struct {
		XMLName xml.Name `xml:"enclosure"`
		URL     string   `xml:"url,attr"`
		Type    string   `xml:"type,attr"`
		Length  uint     `xml:"length,attr"`
	}
)
