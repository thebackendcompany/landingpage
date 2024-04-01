package config

import (
	"github.com/go-batteries/diaper"
	"github.com/rs/zerolog/log"
)

type AppConfig struct {
	Env             string
	Port            int
	LogLevel        string
	UpstashURL      string
	UpstashUserName string
	UpstashPassword string
}

func BuildAppConfig(env string) *AppConfig {
	dc := diaper.DiaperConfig{
		Providers:      diaper.Providers{diaper.EnvProvider{}},
		DefaultEnvFile: "app.env",
	}

	cfgMap, err := dc.ReadFromFile(env, "./config/")
	if err != nil {
		// logrus.WithError(err).Fatal("failed to load config from .env")
		log.Error().Err(err).Msg("failed to read config")
	}

	cfg := &AppConfig{
		Env:             env,
		Port:            cfgMap.MustGetInt("port"),
		LogLevel:        cfgMap.MustGet("log_level").(string),
		UpstashURL:      cfgMap.MustGet("upstash_url").(string),
		UpstashUserName: cfgMap.MustGet("upstash_username").(string),
		UpstashPassword: cfgMap.MustGet("upstash_password").(string),
	}

	return cfg
}
