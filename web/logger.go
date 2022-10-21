package web

import (
	"fmt"
	"github.com/atomAltera/youcaster/logger"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

func newStructuredLogger(log logger.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&structuredLogger{log})
}

type structuredLogger struct {
	log logger.Logger
}

func (l *structuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	logFields := make(map[string]any, 10)

	logFields["ts"] = time.Now().UTC().Format(time.RFC1123)

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields["rid"] = reqID
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	logFields["http_scheme"] = scheme
	logFields["http_proto"] = r.Proto
	logFields["http_method"] = r.Method

	remoteAddr := r.RemoteAddr
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		remoteAddr = ip
	}

	logFields["remote_addr"] = remoteAddr
	logFields["user_agent"] = r.UserAgent()

	logFields["uri"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

	entry := &structuredLoggerEntry{Log: l.log.WithFields(logFields)}
	entry.Log.Infoln("request started")

	return entry
}

type structuredLoggerEntry struct {
	Log logger.Logger
}

func (l *structuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Log = l.Log.WithFields(map[string]any{
		"resp_status": status, "resp_bytes_length": bytes,
		"resp_elapsed_ms": float64(elapsed.Nanoseconds()) / 1000000.0,
	})

	l.Log.Infoln("request complete")
}

func (l *structuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Log = l.Log.WithFields(map[string]any{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}
