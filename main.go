package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/daniel-z-johnson/spotify-playlist-backup/config"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	configPath := flag.String("config", "config.toml", "path to the TOML configuration file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		slog.Error("failed to load configuration", "config_path", *configPath, "error", err)
		os.Exit(1)
	}

	slog.Info(
		"configuration loaded",
		slog.String("config_path", *configPath),
		slog.String("spotify_resource_url", cfg.Spotify.ResourceURL),
	)
}
