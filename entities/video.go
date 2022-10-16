package entities

import "time"

type VideoInfo struct {
	PublishedAt  time.Time
	Title        string
	Description  string
	ThumbnailURL string
}
