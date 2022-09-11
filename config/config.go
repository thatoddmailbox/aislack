package config

import (
	"github.com/BurntSushi/toml"
)

const configFilename = "config.toml"

type Config struct {
	Slack struct {
		AppID          string
		ClientID       string
		ClientSecret   string
		SigningSecret  string
		UserOAuthToken string
	}
	WebhookServer struct {
		Port int
	}
}

var Current Config

func Load() error {
	_, err := toml.DecodeFile(configFilename, &Current)
	if err != nil {
		return err
	}

	return nil
}
