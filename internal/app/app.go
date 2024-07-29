package app

import (
	"context"
	"os"

	"github.com/ShevelevEvgeniy/kafkaManager/config"
	_ "github.com/ShevelevEvgeniy/kafkaManager/docs"
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
		a.log.Error("error occurred on http_server shutting down:", zap.String("error", err.Error()))
		return errors.Wrap(err, "error occurred on http_server shutting down")
	}

	a.log.Info("starting message consumer")

	consumerCtx, consumerCancel := context.WithCancel(ctx)
	err = di.MessageConsumer(consumerCtx).Start(consumerCtx, a.cfg.Kafka)
	if err != nil {
		a.log.Error("error occurred on message consumer shutting down:", zap.String("error", err.Error()))
		consumerCancel()
		return errors.Wrap(err, "error occurred on message consumer shutting down")
	}

	a.log.Info("application started")

	server.Shutdown(ctx, a.log, a.cfg.HTTPServer.StopTimeout)
	consumerCancel()

	return nil
}
