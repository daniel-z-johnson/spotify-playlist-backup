package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/daniel-z-johnson/spotify-playlist-backup/models"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
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

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

func (sc *SpotifyController) AuthorizationURL(w http.ResponseWriter, r *http.Request) {
	authURL, err := sc.SpotifyService.GetAuthorizationURL()
	if err != nil {
		sc.Logger.Error("Failed to get Spotify authorization URL", slog.String("error", err.Error()))
		http.Error(w, "Failed to get authorization URL", http.StatusInternalServerError)
		return
	}
	sc.Templates.Link.Execute(w, r, SpotifyOAuthLink{
		OAuthLink: authURL,
	})
}

func (sc *SpotifyController) GetToken(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	sc.Logger.Info("State received", slog.String("state", state))
	tokenRequest := url.Values{}
	tokenRequest.Add("code", code)
	tokenRequest.Add("grant_type", "authorization_code")
	tokenRequest.Add("redirect_uri", sc.SpotifyService.Configuration.SpotifyConfig.RedirectURI)
	basicPlain := fmt.Sprintf("%s:%s", sc.SpotifyService.Configuration.SpotifyConfig.ClientID, sc.SpotifyService.Configuration.SpotifyConfig.ClientSecret)
	basicBase64 := base64.StdEncoding.EncodeToString([]byte(basicPlain))
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(tokenRequest.Encode()))
	if err != nil {
		sc.Logger.Error("Failed to parse Spotify token URL", slog.String("error", err.Error()))
		http.Error(w, "Failed to process request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", basicBase64))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		sc.Logger.Error("Failed to get token from Spotify", slog.String("error", err.Error()))
		http.Error(w, "Failed to get token", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		sc.Logger.Error("Didn't get a 200 OK resp from spotify", slog.String("status", resp.Status))
		http.Error(w, "Failed to get token", http.StatusInternalServerError)
		return
	}
	tokenResponse := &TokenResponse{}
	err = json.NewDecoder(resp.Body).Decode(tokenResponse)
	if err != nil {
		sc.Logger.Error("Failed to decode token response", slog.String("error", err.Error()))
		http.Error(w, "Failed to decode token response", http.StatusInternalServerError)
		return
	}
	sc.Logger.Info("Successfully received token")
	fmt.Fprintf(w, "Retrieved token")
}
