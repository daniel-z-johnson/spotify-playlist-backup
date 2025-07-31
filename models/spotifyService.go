package models

import (
	spotConfig "github.com/daniel-z-johnson/spotify-playlist-backup/config"
	"log/slog"
	"net/url"
)

type SpotifyService struct {
	Logger        *slog.Logger
	Configuration *spotConfig.AppConfiguration
}

func (s *SpotifyService) GetAuthorizationURL() (string, error) {
	authURL, err := url.Parse("https://accounts.spotify.com/authorize")
	if err != nil {
		s.Logger.Error("Failed to parse Spotify authorization URL", slog.String("error", err.Error()))
		return "", err
	}
	params := url.Values{}
	params.Add("client_id", s.Configuration.SpotifyConfig.ClientID)
	params.Add("response_type", "code")
	params.Add("redirect_uri", s.Configuration.SpotifyConfig.RedirectURI)
	params.Add("scope", "playlist-read-private playlist-read-collaborative")
	params.Add("state", "not-implemented")
	authURL.RawQuery = params.Encode()
	return authURL.String(), nil
}
