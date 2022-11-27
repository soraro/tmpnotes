package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type specification struct {
	Port               int    `default:"5000" envconfig:"PORT"`
	RedisUrl           string `default:"redis://localhost:6379" envconfig:"REDIS_URL"`
	EnableHsts         bool   `split_words:"true"`
	MaxLength          int    `default:"1000" split_words:"true"`
	UiMaxLength        int    `default:"512" split_words:"true"`
	MaxExpire          int    `default:"24" split_words:"true"`
	SlackToken         string `split_words:"true"`
	SlackSigningSecret string `split_words:"true"`
}

var Config specification

func GetConfig() error {
	err := envconfig.Process("tmpnotes", &Config)
	if err != nil {
		return err
	}

	if Config.UiMaxLength > Config.MaxLength {
		return fmt.Errorf("UiMaxLength %v should not be greater than MaxLength %v", Config.UiMaxLength, Config.MaxLength)
	}
	return nil
}
