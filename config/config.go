package config

import (
	"encoding/json"
	"os"
)

type AppConfiguration struct {
	SpotifyConfig struct {
		ClientID     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
		RedirectURI  string `json:"redirectUri"`
	} `json:"spotify"`
}

func LoadConfig(fileLoc string) (*AppConfiguration, error) {
	config := &AppConfiguration{}
	file, err := os.Open(fileLoc)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
