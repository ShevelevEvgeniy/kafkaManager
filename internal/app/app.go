package app

import (
	"context"

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
	di := NewDiContainer(a.cfg, a.log)

	router := initRouter(ctx, di)

	server := NewServer(a.cfg, router)
	err := server.Run(a.log, a.cfg)
	if err != nil {
		a.log.Error("error occurred on http_server shutting down:", zap.String("error", err.Error()))
		return errors.Wrap(err, "error occurred on http_server shutting down")
	}

	a.log.Info("starting message consumer")

	a.registeredEvents(ctx, di)

	consumerCtx, consumerCancel := context.WithCancel(ctx)
	di.MessageConsumer(consumerCtx).Start(consumerCtx)

	a.log.Info("application started")

	server.Shutdown(ctx, a.log, a.cfg.HTTPServer.StopTimeout)
	consumerCancel()

	a.log.Info("application stopped")

	return nil
}

func (a *App) registeredEvents(ctx context.Context, di DiContainer) {
	a.log.Info("registered events")

	eventDispatcher := di.EventDispatcher(ctx)
	eventDispatcher.Subscribe("order_status", di.MessageStatusHandler(ctx))
}
