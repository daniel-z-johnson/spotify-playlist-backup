// Package config loads application configuration from TOML files.
package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// Config contains all application configuration.
type Config struct {
	Spotify SpotifyConfig `toml:"spotify"`
}

// SpotifyConfig contains the settings used to access Spotify.
type SpotifyConfig struct {
	ClientID     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
	RedirectURL  string `toml:"redirect_url"`
	AuthURL      string `toml:"auth_url"`
	ResourceURL  string `toml:"resource_url"`
}

// Load reads and parses a TOML configuration file at path.
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config %q: %w", path, err)
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config %q: %w", path, err)
	}

	return cfg, nil
}
