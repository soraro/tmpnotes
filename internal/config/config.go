package config

import (
	"github.com/kelseyhightower/envconfig"
)

type specification struct {
	Port     int    `default:"5000"`
	RedisUrl string `default:"redis://localhost:6379" envconfig:"REDIS_URL"`
}

var Config specification

func GetConfig() error {
	return envconfig.Process("tmpnotes", &Config)
}
