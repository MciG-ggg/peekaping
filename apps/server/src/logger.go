package main

import (
	"context"
	"peekaping/src/config"
	"time"

	zaploki "github.com/paul-milne/zap-loki"
	"go.uber.org/zap"
)

func ProvideLogger(cfg *config.Config) (*zap.SugaredLogger, error) {
	zapConfig := zap.NewProductionConfig()
	var logger *zap.Logger
	var err error

	if cfg.DevMode {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zapConfig.Build()
	}
	if err != nil {
		return nil, err
	}

	if cfg.LokiURL != "" {
		loki := zaploki.New(context.Background(), zaploki.Config{
			Url:          cfg.LokiURL,
			BatchMaxSize: 1000,
			BatchMaxWait: 10 * time.Second,
			Labels:       map[string]string{"service_name": "peekaping"},
		})

		logger, err = loki.WithCreateLogger(zapConfig)
		if err != nil {
			return nil, err
		}
	}

	return logger.Sugar(), nil
}
