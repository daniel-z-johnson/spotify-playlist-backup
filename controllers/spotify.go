package controllers

import (
	"github.com/daniel-z-johnson/spotify-playlist-backup/models"
	"log/slog"
	"net/http"
)

type SpotifyOAuthLink struct {
	OAuthLink string
}
type SpotifyController struct {
	Templates struct {
		Link Template
	}
	SpotifyService *models.SpotifyService
	Logger         *slog.Logger
}

func (sc *SpotifyController) AuthorizationURL(w http.ResponseWriter, r *http.Request) {
	sc.Templates.Link.Execute(w, r, SpotifyOAuthLink{
		OAuthLink: sc.SpotifyService.GetAuthorizationURL(),
	})
}
