package config

import "testing"

func TestLoadConfig(t *testing.T) {
	conf, err := LoadConfig("example-config.json")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	empty := AppConfiguration{}
	if conf == nil {
		t.Fatalf("configuration should not be nil")
	}
	if conf.SpotifyConfig == empty.SpotifyConfig {
		t.Fatalf("SpotifyConfig should not be empty")
	}
}
