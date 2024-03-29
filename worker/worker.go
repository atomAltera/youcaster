package worker

import (
	"context"
	e "github.com/atomAltera/youcaster/entities"
	"github.com/atomAltera/youcaster/logger"
	"math"
	"time"
)

const maxRetries = 3

type Worker struct {
	log        logger.Logger
	storage    RequestsStore
	infoGetter InfoGetter
	downloader Downloader

	failedRequestsChan chan e.FailedRequest
}

func NewWorker(l logger.Logger, s RequestsStore, ig InfoGetter, dl Downloader) *Worker {
	return &Worker{
		log:                l,
		storage:            s,
		infoGetter:         ig,
		downloader:         dl,
		failedRequestsChan: make(chan e.FailedRequest, 100),
	}
}

func (w *Worker) StartListenRequests(rc <-chan e.Request) {
	go func() {
		for r := range rc {
			w.processIncomingRequest(r)
		}
	}()
}

func (w *Worker) StartProcessingRequests() <-chan e.FailedRequest {
	go w.processingLoop()
	return w.failedRequestsChan
}

func (w *Worker) processingLoop() {
	for {
		w.processingCycle()
		time.Sleep(10 * time.Second)
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
	if r.Attempts > 0 {
		wait := time.Duration(math.Pow(5, float64(r.Attempts))) * time.Second
		if time.Since(r.UpdatedAt) < wait {
			return
		}
	}

	if r.Status == e.RequestStatusNew {
		r.Status = e.RequestStatusDownloading
		if ok := w.updateRequest(r); !ok {
			return
		}
	}

	l := w.log.WithFields(map[string]any{
		"attempt":          r.Attempts + 1,
		"id":               r.ID,
		"youtube_video_id": r.YoutubeVideoID,
	})

	l.Info("getting video info")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	info, err := w.infoGetter.GetInfo(ctx, r.YoutubeVideoID)
	if err != nil {
		l.WithError(err).Error("failed to get youtube video info")
		w.requestFiled(r, err)

		r.Attempts += 1
		r.Status = e.RequestStatusNew
		if r.Attempts >= maxRetries {
			r.Status = e.RequestStatusFailed
		}

		w.updateRequest(r)
		return
	}

	r.VideoInfo = info
	l.SetField("duration", info.Duration)

	if ok := w.updateRequest(r); !ok {
		return
	}

	r.FileName = w.generateFileName(r)

	l.Info("downloading video")
	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	fileSize, err := w.downloader.Download(ctx, r.YoutubeVideoID, r.FileName)
	if err != nil {
		r.Attempts += 1

		l.WithError(err).Error("failed to download video")

		r.Status = e.RequestStatusNew
		r.LastAttemptAt = time.Now()

		if r.Attempts > maxRetries {
			r.Status = e.RequestStatusFailed
			w.requestFiled(r, err)
		}

		w.updateRequest(r)
		return
	}

	r.FileSize = fileSize
	r.Status = e.RequestStatusDone

	l.SetField("file_size", fileSize)

	if ok := w.updateRequest(r); !ok {
		return
	}

	l.Info("video downloaded")
}

func (w *Worker) updateRequest(r e.Request) bool {
	r.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	if err := w.storage.Update(ctx, r); err != nil {
		w.log.WithError(err).Error("failed to update request")
		w.requestFiled(r, err)
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
		w.requestFiled(r, err)
	}
}

func (w *Worker) generateID(r e.Request) string {
	return "ytb_" + r.YoutubeVideoID
}

func (w *Worker) generateFileName(r e.Request) string {
	return "ytb_" + r.YoutubeVideoID + ".mp3"
}

func (w *Worker) requestFiled(r e.Request, err error) {
	select {
	case w.failedRequestsChan <- e.FailedRequest{
		Request: r,
		Error:   err,
	}:
		break
	default:
		w.log.Error("failed requests chan is full")
	}
}
