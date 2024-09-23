package main

import (
	"github.com/mdshahjahanmiah/explore-go/di"
	eHttp "github.com/mdshahjahanmiah/explore-go/http"
	"github.com/mdshahjahanmiah/explore-go/logging"
	"github.com/mdshahjahanmiah/sales-manager-scheduler/pkg/config"
	"github.com/mdshahjahanmiah/sales-manager-scheduler/pkg/db"
	sales_managers "github.com/mdshahjahanmiah/sales-manager-scheduler/pkg/sales-managers"
	"go.uber.org/dig"
	"log/slog"
)

func main() {
	c := di.New()

	c.Provide(func() (config.Config, error) {
		conf, err := config.Load()
		if err != nil {
			slog.Error("failed to load configuration", "err", err)
			return config.Config{}, err
		}
		return conf, nil
	})

	slog.Info("application configuration is loaded successfully")

	c.Provide(func(conf config.Config) (*logging.Logger, error) {
		logger, err := logging.NewLogger(conf.LoggerConfig)
		if err != nil {
			slog.Error("initializing logger", "err", err)
			return nil, err
		}

		return logger, nil
	})

	slog.Info("application logger is initialized successfully")

	c.Provide(func(conf config.Config, logger *logging.Logger) (*db.DB, error) {
		db, err := db.NewDB(conf.PostgresDSN, logger)
		if err != nil {
			logger.Error("database initialization", "err", err.Error())
			return nil, err
		}
		return db, nil
	})

	c.Provide(func(config config.Config) *eHttp.ServerConfig {
		return &eHttp.ServerConfig{
			HttpAddress: config.HttpAddress,
		}
	})

	c.Provide(func(config config.Config, logger *logging.Logger, db *db.DB) (sales_managers.Service, error) {
		service, err := sales_managers.NewService(config, logger, db)
		if err != nil {
			slog.Error("initializing sales-managers service", "err", err)
			return nil, err
		}
		return service, nil
	})

	c.ProvideMonitoringEndpoints("endpoint")

	c.Provide(sales_managers.MakeHandler, dig.Group("endpoint"))

	c.Invoke(func(in struct {
		dig.In
		Conf         config.Config
		ServerConfig *eHttp.ServerConfig
		Endpoints    []eHttp.Endpoint `group:"endpoint"`
	}) {
		server := eHttp.NewServer(in.ServerConfig, in.Endpoints, nil)
		c.Provide(func() di.StartCloser { return server }, dig.Group("startclose"))
	})

	c.Start()

}
