package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Config struct {
	Mode       ProjectMode
	HTTPServer HTTPServer
	DB         DB
	Kafka      Kafka
}

func MustLoad(log *zap.Logger) (*Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, errors.Wrap(err, "failed to load config")
	}

	logNonSecretConfig(log, &cfg)

	return &cfg, nil
}

func logNonSecretConfig(log *zap.Logger, cfg *Config) {
	log.Info("Initialized config",
		zap.String("Mode", cfg.Mode.Env),
		zap.String("HTTPServer Port", cfg.HTTPServer.Port),
		zap.Duration("HTTPServer Timeout", cfg.HTTPServer.Timeout),
		zap.Duration("HTTPServer IdleTimeout", cfg.HTTPServer.IdleTimeout),
		zap.Duration("HTTPServer StopTimeout", cfg.HTTPServer.StopTimeout),
		zap.String("DB Host", cfg.DB.Host),
		zap.String("DB Port", cfg.DB.Port),
		zap.String("DB Name", cfg.DB.DBName),
		zap.String("DB DriverName", cfg.DB.DriverName),
		zap.String("Kafka Broker", cfg.Kafka.Broker),
		zap.String("Kafka Topic", cfg.Kafka.Topic),
	)
}
