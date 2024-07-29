package main

import (
	"context"
	"os"

	"github.com/ShevelevEvgeniy/kafkaManager/config"
	_ "github.com/ShevelevEvgeniy/kafkaManager/docs"
	"github.com/ShevelevEvgeniy/kafkaManager/internal/app"
	"github.com/ShevelevEvgeniy/kafkaManager/lib/logger/uber_zap"
	"go.uber.org/zap"
)

// @title Kafka Manager API
// @version 1.0
// @description This API manages Kafka messages by allowing you to create orders and check their status.
// @contact.name API Support
// @contact.email Z_shevelev@mail.ru
// @BasePath /api/v1
func main() {
	ctx := context.Background()

	log, stop := uber_zap.InitLogger(os.Getenv("ENV_TYPE"))
	defer stop()

	log.Info("Initialized logger")

	cfg, err := config.MustLoad(log)
	if err != nil {
		log.Fatal("failed to load config", zap.Error(err))
		os.Exit(1)
	}

	log.Info("Started application")

	application := app.NewApp(cfg, log)
	if err = application.Run(ctx); err != nil {
		panic(err)
	}
}
