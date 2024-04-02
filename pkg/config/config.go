package config

import (
	"thebackendcompany/pkg/creepto"

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

	GoogleCreds      string
	EmailLeadsDbName string
	DomainName       string
}

func BuildAppConfig(env string) *AppConfig {
	dc := diaper.DiaperConfig{
		Providers:      diaper.Providers{diaper.EnvProvider{}},
		DefaultEnvFile: "app.env",
	}

	cfgMap, err := dc.ReadFromFile(env, "./config/")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read config")
	}

	cfg := &AppConfig{
		Env:              env,
		Port:             cfgMap.MustGetInt("port"),
		LogLevel:         cfgMap.MustGet("log_level").(string),
		UpstashURL:       cfgMap.MustGet("upstash_url").(string),
		UpstashUserName:  cfgMap.MustGet("upstash_username").(string),
		UpstashPassword:  cfgMap.MustGet("upstash_password").(string),
		EmailLeadsDbName: cfgMap.MustGet("email_leads_db_name").(string),
		DomainName:       cfgMap.MustGet("domain_name").(string),
	}

	googleCredsFile := cfgMap.MustGet("google_creds_file").(string)
	masterKey := cfgMap.MustGet("master_key").(string)

	plainText, err := creepto.Decrypt(googleCredsFile, masterKey)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to decrypt creds file")
	}

	// add a validation to validate authenticity of file
	cfg.GoogleCreds = plainText

	return cfg
}
