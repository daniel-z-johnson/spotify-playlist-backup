package server

import (
	"log/slog"
	"net/http"

	"github.com/daniel-z-johnson/spotify-playlist-backup/config"
	"github.com/go-chi/chi/v5"
)

func route(log *slog.Logger, conf *config.Config) http.Handler {
	log.Info("setting up routes")
	r := chi.NewRouter()

	return r
}
