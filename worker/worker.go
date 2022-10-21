package worker

import (
	"context"
	e "github.com/atomAltera/youcaster/entities"
	"github.com/atomAltera/youcaster/logger"
	"time"
)

const maxRetries = 5

type Worker struct {
	log        logger.Logger
	storage    RequestsStore
	infoGetter InfoGetter
	downloader Downloader
}

func NewWorker(l logger.Logger, s RequestsStore, ig InfoGetter, dl Downloader) *Worker {
	return &Worker{
		log:        l,
		storage:    s,
		infoGetter: ig,
		downloader: dl,
	}
}

func (w *Worker) StartListenRequests(rc <-chan e.Request) {
	go func() {
		for r := range rc {
			w.processIncomingRequest(r)
		}
	}()
}

func (w *Worker) StartProcessingRequests() {
	go w.processingLoop()
}

func (w *Worker) processingLoop() {
	for {
		w.processingCycle()
		time.Sleep(1 * time.Second)
	}
}

func (w *Worker) processingCycle() {
	ss := []e.RequestStatus{
		e.RequestStatusNew,
		e.RequestStatusDownloading,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	rs, err := w.storage.List(ctx, ss)
	if err != nil {
		w.log.WithError(err).Error("failed to fetch requests")
		return
	}

	for _, r := range rs {
		w.processRequest(r)
	}
}

func (w *Worker) processRequest(r e.Request) {
	if r.Status == e.RequestStatusNew {
		r.Status = e.RequestStatusDownloading
		if ok := w.updateRequest(r); !ok {
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	info, err := w.infoGetter.GetInfo(ctx, r.YoutubeVideoID)
	if err != nil {
		w.log.WithError(err).Error("failed to get youtube video info")

		r.Attempts += 1
		r.Status = e.RequestStatusNew
		if r.Attempts > maxRetries {
			r.Status = e.RequestStatusFailed
		}

		w.updateRequest(r)
		return
	}

	r.VideoInfo = info

	if ok := w.updateRequest(r); !ok {
		return
	}

	r.FileName = w.generateFileName(r)

	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	ytUrl := "https://www.youtube.com/watch?v=" + r.YoutubeVideoID
	fileSize, err := w.downloader.Download(ctx, ytUrl, r.FileName)
	if err != nil {
		w.log.WithError(err).Error("failed to download video")

		r.Attempts += 1
		r.Status = e.RequestStatusNew
		if r.Attempts > maxRetries {
			r.Status = e.RequestStatusFailed
		}

		w.updateRequest(r)
		return
	}

	r.FileSize = fileSize
	r.Status = e.RequestStatusDone

	w.updateRequest(r)
}

func (w *Worker) updateRequest(r e.Request) bool {
	r.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	if err := w.storage.Update(ctx, r); err != nil {
		w.log.WithError(err).Error("failed to update request")
		return false
	}

	return true
}

func (w *Worker) processIncomingRequest(r e.Request) {
	now := time.Now()

	r.ID = w.generateID(r)
	r.CreatedAt = now
	r.UpdatedAt = now
	r.Status = e.RequestStatusNew

	l := w.log.WithFields(map[string]any{
		"id":               r.ID,
		"youtube_video_id": r.YoutubeVideoID,
		"tg_chat_id":       r.TgChatID,
		"tg_message_id":    r.TgMessageID,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	l.Info("saving request")
	err := w.storage.Create(ctx, r)
	if err != nil {
		l.WithError(err).Error("failed to save requests")
	}
}

func (w *Worker) generateID(r e.Request) string {
	return "ytb_" + r.YoutubeVideoID
}

func (w *Worker) generateFileName(r e.Request) string {
	return "ytb_" + r.YoutubeVideoID + ".mp3"
}
