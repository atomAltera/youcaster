package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	e "github.com/atomAltera/youcaster/entities"
	"net/http"
	"net/url"
	"time"
)

type InfoGetter struct {
	apiKey string
	client HttpClient
}

func NewInfoGetter(apiKey string, client HttpClient) (*InfoGetter, error) {
	if client == nil {
		client = http.DefaultClient
	}

	return &InfoGetter{
		apiKey: apiKey,
		client: client,
	}, nil
}

func (i *InfoGetter) GetInfo(ctx context.Context, id string) (*e.VideoInfo, error) {
	pu, err := url.Parse("https://www.googleapis.com/youtube/v3/videos")
	if err != nil {
		return nil, err
	}

	q := pu.Query()
	q.Set("key", i.apiKey)
	q.Set("part", "snippet")
	q.Set("id", url.QueryEscape(id))
	pu.RawQuery = q.Encode()
	u := pu.String()

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	res, err := i.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed send http request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http request failed with status code: %d (%v)", res.StatusCode, res.Status)
	}

	var payload VideoListResponse
	if err = json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	if len(payload.Items) == 0 {
		return nil, fmt.Errorf("video not found")
	}

	item := payload.Items[0]

	publishedAt, err := time.Parse("2006-01-02T15:04:05Z", item.Snippet.PublishedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse published at: %w", err)
	}

	thumbnailURL := ""
	if standardThumbnail, ok := item.Snippet.Thumbnails["standard"]; ok {
		thumbnailURL = standardThumbnail.URL
	}

	return &e.VideoInfo{
		PublishedAt:  publishedAt,
		Title:        item.Snippet.Title,
		Description:  item.Snippet.Description,
		ThumbnailURL: thumbnailURL,
	}, nil
}
