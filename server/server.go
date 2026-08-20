package server

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/daniel-z-johnson/spotify-playlist-backup/config"
)

type Server struct {
	log  *slog.Logger
	http *http.Server
}

func New(log *slog.Logger, conf *config.Config) *Server {
	address := conf.Address
	if strings.TrimSpace(address) == "" {
		address = ":1117"
	}
	handler := route(log, conf)
	return &Server{
		log: log,
		http: &http.Server{
			Addr:              address,
			Handler:           handler,
			ReadHeaderTimeout: 5 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	s.log.Info("starting server", "address", s.http.Addr)
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
