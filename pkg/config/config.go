package config

import (
	"flag"
	"github.com/mdshahjahanmiah/explore-go/logging"
)

type Config struct {
	HttpAddress  string
	PostgresDSN  string
	LoggerConfig logging.LoggerConfig
}

func Load() (Config, error) {
	fs := flag.NewFlagSet("", flag.ExitOnError)

	httpAddress := fs.String("http.public.address", "localhost:3000", "HTTP listen address for all specified endpoints.")
	postgresDSN := fs.String("dsn", "postgres://postgress:mypassword123!@localhost:5432/coding-challenge?sslmode=disable", "DB address")

	loggerConfig := logging.LoggerConfig{}
	fs.StringVar(&loggerConfig.CommandHandler, "logger.handler.type", "json", "handler type e.g json, otherwise default will be text type")
	fs.StringVar(&loggerConfig.LogLevel, "logger.log.level", "debug", "log level wise logging with fatal log")

	config := Config{
		HttpAddress:  *httpAddress,
		PostgresDSN:  *postgresDSN,
		LoggerConfig: loggerConfig,
	}

	return config, nil
}
