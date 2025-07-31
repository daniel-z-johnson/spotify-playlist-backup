package main

import (
	"github.com/daniel-z-johnson/spotify-playlist-backup/controllers"
	"github.com/daniel-z-johnson/spotify-playlist-backup/models"
	"github.com/daniel-z-johnson/spotify-playlist-backup/templates"
	"github.com/daniel-z-johnson/spotify-playlist-backup/views"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("Start")

	spotifyController := &controllers.SpotifyController{
		Logger:         logger,
		SpotifyService: &models.SpotifyService{},
	}
	spotifyController.Templates.Link = views.Must(views.ParseFS(templates.TemplatesFS, logger, "main-layout.gohtml", "oauth-link.gohtml"))

	logger.Info("Server starting on port 1117")

	r := chi.NewRouter()
	r.Get("/", spotifyController.AuthorizationURL)
	http.ListenAndServe(":1117", r)
}
