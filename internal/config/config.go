package config

import (
	"github.com/kelseyhightower/envconfig"
)

type specification struct {
	Port       int    `default:"5000" envconfig:"PORT"`
	RedisUrl   string `default:"redis://localhost:6379" envconfig:"REDIS_URL"`
	EnableHsts bool   `split_words:"true"`
}

var Config specification

func GetConfig() error {
	return envconfig.Process("tmpnotes", &Config)
}
