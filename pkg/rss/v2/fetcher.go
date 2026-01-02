package v2

import "net/http"

type Fetcher struct {
	client *http.Client
}

func NewFetcher(client *http.Client) *Fetcher {
	return &Fetcher{
		client: client,
	}
}

func (f *Fetcher) Fetch(URL string) (*Feed, error) {
	resp, err := f.client.Get(URL)

	if err != nil {
		return nil, err
	}

	feed, err := Parse(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	return feed, nil
}
