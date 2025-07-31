package models

type SpotifyService struct {
}

func (s *SpotifyService) GetAuthorizationURL() string {
	return "https://accounts.spotify.com/authorize"
}
