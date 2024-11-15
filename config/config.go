package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Env             string
	ParentAPIURL    string
	Port            string
	CacheExpiration time.Duration
}

func Init() *Config {
	// Load .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Debug().Msg("No .env file found")
	}

	// Set default values
	cfg := &Config{
		Env:             getEnvWithDefault("ENV", "prod"),
		ParentAPIURL:    "https://bsky-search.jazco.io/stats",
		Port:            getEnvWithDefault("PORT", "8080"),
		CacheExpiration: time.Second * 60,
	}

	// Configure logging
	setupLogging(cfg.Env)

	return cfg
}

func setupLogging(env string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Set up console writer
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	// Set global logger
	log.Logger = log.Output(consoleWriter)

	// Set log level based on environment
	if env == "dev" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Debug logging enabled")
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		log.Info().Str("env", env).Msg("Production logging enabled")
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
