package app

import (
	"context"
	"os"

	"github.com/ShevelevEvgeniy/kafkaManager/config"
	dbConn "github.com/ShevelevEvgeniy/kafkaManager/internal/postgres/db_connection"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type DiContainer struct {
	cfg *config.Config
	log *zap.Logger
	db  *pgxpool.Pool
}

func NewDiContainer(cfg *config.Config, log *zap.Logger) DiContainer {
	return DiContainer{
		cfg: cfg,
		log: log,
	}
}

func (di *DiContainer) DB(ctx context.Context) *pgxpool.Pool {
	if di.db == nil {
		db, err := dbConn.Connect(ctx, di.cfg.DB)
		if err != nil {
			di.log.Fatal("failed to connect to db", zap.Error(err))
			os.Exit(1)
		}

		di.log.Info("connected to db", zap.String("database", di.cfg.DB.DBName))
		di.db = db
	}

	return di.db
}
