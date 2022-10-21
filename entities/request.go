package entities

import "time"

type RequestStatus = string

const (
	RequestStatusNew         RequestStatus = "new"
	RequestStatusDownloading RequestStatus = "downloading"
	RequestStatusFailed      RequestStatus = "failed"
	RequestStatusDone        RequestStatus = "done"
)

type Request struct {
	ID     string        `json:"id" bson:"id"`
	Status RequestStatus `json:"status" bson:"status"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`

	YoutubeVideoID string `json:"youtube_video_id" bson:"youtube_video_id"`

	TgChatID    int64 `json:"tg_chat_id" bson:"tg_chat_id"`
	TgMessageID int   `json:"tg_message_id" bson:"tg_message_id"`

	Attempts int `json:"attempts" bson:"attempts"`

	VideoInfo *VideoInfo `json:"video_info" bson:"video_info"`
	FileName  string     `json:"file_name" bson:"file_name"`
	FileSize  int64      `json:"file_size" bson:"file_size"`
}
