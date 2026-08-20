package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/daniel-z-johnson/spotify-playlist-backup/config"
)

func TestLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	contents := `[spotify]
client_id = "client-id"
client_secret = "client-secret"
redirect_url = "http://127.0.0.1:8080/callback"
auth_url = "https://accounts.spotify.com"
resource_url = "https://api.spotify.com"
[database]
location = "test.db"
`
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("write test config: %v", err)
	}

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	want := config.SpotifyConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		RedirectURL:  "http://127.0.0.1:8080/callback",
		AuthURL:      "https://accounts.spotify.com",
		ResourceURL:  "https://api.spotify.com",
	}
	if cfg.Spotify != want {
		t.Errorf("Load().Spotify = %#v, want %#v", cfg.Spotify, want)
	}
	expectedDBLocation := "test.db"
	if cfg.Database.Location != expectedDBLocation {
		t.Errorf("Load().Database.Location = %q, want %q", cfg.Database.Location, expectedDBLocation)
	}
}

func TestLoadMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.toml")

	_, err := config.Load(path)
	if err == nil {
		t.Fatal("Load() error = nil, want an error")
	}
	if !strings.Contains(err.Error(), path) {
		t.Errorf("Load() error = %q, want it to contain %q", err, path)
	}
}

func TestLoadMalformedTOML(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte("[spotify\n"), 0o600); err != nil {
		t.Fatalf("write test config: %v", err)
	}

	_, err := config.Load(path)
	if err == nil {
		t.Fatal("Load() error = nil, want an error")
	}
	if !strings.Contains(err.Error(), path) {
		t.Errorf("Load() error = %q, want it to contain %q", err, path)
	}
}
