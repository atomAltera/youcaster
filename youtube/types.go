package youtube

import "net/http"

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type VideoListResponse struct {
	Items []Item `json:"items"`
}

type Item struct {
	ID      string  `json:"id"`
	Snippet Snippet `json:"snippet"`
}

type Snippet struct {
	Title       string               `json:"title"`
	PublishedAt string               `json:"publishedAt"`
	Description string               `json:"description"`
	Thumbnails  map[string]Thumbnail `json:"thumbnails"`
}

type Thumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
