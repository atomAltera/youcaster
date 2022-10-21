package entities

import "time"

type VideoInfo struct {
	PublishedAt  time.Time     `json:"published_at" bson:"published_at"`
	Title        string        `json:"title" bson:"title"`
	Description  string        `json:"description" bson:"description"`
	ThumbnailURL string        `json:"thumbnail_url" bson:"thumbnail_url"`
	Duration     time.Duration `json:"duration" bson:"duration"`
}
