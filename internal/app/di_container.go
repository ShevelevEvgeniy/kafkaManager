package app

import (
	"github.com/ShevelevEvgeniy/kafkaManager/config"
	"go.uber.org/zap"
)

type DiContainer struct {
	cfg *config.Config
	log *zap.Logger
}

func NewDiContainer(cfg *config.Config, log *zap.Logger) DiContainer {
	return DiContainer{
		cfg: cfg,
		log: log,
	}
}
