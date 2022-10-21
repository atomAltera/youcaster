package web

import (
	_ "embed"
	e "github.com/atomAltera/youcaster/entities"
	"github.com/atomAltera/youcaster/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"net/http"
)

//go:embed static/logo.png
var logo []byte

type Server struct {
	l          logger.Logger
	r          chi.Router
	store      RequestsStore
	fileReader FileReader
	builder    FeedBuilder
}

func NewServer(
	l logger.Logger,
	store RequestsStore,
	fr FileReader,
	fb FeedBuilder,
) *Server {
	s := &Server{
		l:          l,
		store:      store,
		fileReader: fr,
		builder:    fb,
		r:          chi.NewRouter(),
	}

	// Middleware
	s.r.Use(newStructuredLogger(l))
	s.r.Use(middleware.Recoverer)

	// Serve feed
	s.r.Get("/logo.png", s.logoHandler)
	s.r.Get("/feed", s.feedHandler)
	s.r.Get("/files/{filename}", s.fileHandler)

	return s
}

// Listen starts the server on address addr.
func (s *Server) Listen(addr string) error {
	return http.ListenAndServe(addr, s.r)
}

func (s *Server) logoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", string(len(logo)))
	_, _ = w.Write(logo)
}

func (s *Server) feedHandler(w http.ResponseWriter, r *http.Request) {
	rs, err := s.store.List(r.Context(), []e.RequestStatus{e.RequestStatusDone})
	if err != nil {
		s.l.WithError(err).Error("failed to fetch requests")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	xml, err := s.builder.BuildFeed(rs)
	if err != nil {
		s.l.WithError(err).Error("failed to build feed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")

	_, err = io.WriteString(w, xml)
	if err != nil {
		s.l.WithError(err).Error("Error writing feed")
	}
}

func (s *Server) fileHandler(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	file, err := s.fileReader.ReadFile(r.Context(), filename)
	if err != nil {
		s.l.WithError(err).Error("failed to read file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "audio/mpeg")

	_, err = io.Copy(w, file)
	if err != nil {
		s.l.WithError(err).Error("Error writing file")
	}
}
