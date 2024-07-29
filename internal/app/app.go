package app

import (
	"context"
	"os"

	"github.com/ShevelevEvgeniy/kafkaManager/config"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type App struct {
	cfg *config.Config
	log *zap.Logger
}

func NewApp(cfg *config.Config, log *zap.Logger) *App {
	return &App{
		cfg: cfg,
		log: log,
	}
}

func (a *App) Run(ctx context.Context) error {
	certFile := os.Getenv("HTTP_SERVER_CERT_FILE")
	keyFile := os.Getenv("HTTP_SERVER_KEY_FILE")

	di := NewDiContainer(a.cfg, a.log)

	router := initRouter(ctx, di)

	server := NewServer(a.cfg, router)
	err := server.Run(a.log, a.cfg, certFile, keyFile)
	if err != nil {
		a.log.Error("error occurred on server shutting down:", zap.String("error", err.Error()))
		return errors.Wrap(err, "error occurred on server shutting down")
	}

	a.log.Info("application started")

	server.Shutdown(ctx, a.log, a.cfg.HTTPServer.StopTimeout)

	return nil
}
