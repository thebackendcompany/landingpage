package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"strings"
	"time"

	"thebackendcompany/pkg/config"

	"thebackendcompany/app/core/events"

	"github.com/rs/zerolog/log"
)

func main() {
	env := os.Getenv("ENVIRONMENT")
	env = strings.ToLower(env)

	if env == "" {
		env = "local"
	}

	cfg := config.BuildAppConfig(env)

	interval := flag.Int64("interval", 15, "interval in seconds")
	flag.Parse()

	mq, err := events.NewKafkaEventConsumer(
		&events.KafkaConfig{
			Servers:      cfg.UpstashURL,
			Topic:        "events-hook", // take from config
			SaslUserName: cfg.UpstashUserName,
			SaslPassword: cfg.UpstashPassword,
		},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to upstash kafka")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	consumerData, err := mq.Consume(
		ctx,
		time.Duration(*interval)*time.Second,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize consumer")
	}

	log.Info().Int64("interval", *interval).Msg("consuming")

	for result := range consumerData {
		if result.Err != nil {
			log.Error().Err(result.Err).Msg("failed to consume from producer")

			continue
		}

		log.Info().Str("received ", string(result.Message)).Msg("")
	}

	// Start server
	// Wait for interrupt signal to gracefully shutdown the server
	<-ctx.Done()
}
