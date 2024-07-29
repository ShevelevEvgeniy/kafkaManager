package app

import (
	"context"
	"os"

	"github.com/ShevelevEvgeniy/kafkaManager/config"
	"github.com/ShevelevEvgeniy/kafkaManager/internal/clients/kafka"
	"github.com/ShevelevEvgeniy/kafkaManager/internal/http_server/api/v1/handlers"
	"github.com/ShevelevEvgeniy/kafkaManager/internal/http_server/events"
	mesTrackerRepo "github.com/ShevelevEvgeniy/kafkaManager/internal/postgres/repository/message_tracker_repository"
	repoInterfaces "github.com/ShevelevEvgeniy/kafkaManager/internal/postgres/repository/repository_interfaces"
	"github.com/ShevelevEvgeniy/kafkaManager/internal/service/order_service"
	servInterfaces "github.com/ShevelevEvgeniy/kafkaManager/internal/service/service_interfaces"
	"github.com/ShevelevEvgeniy/kafkaManager/internal/service/statuses_service"
	"github.com/go-playground/validator/v10"

	dbConn "github.com/ShevelevEvgeniy/kafkaManager/internal/postgres/db_connection"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type DiContainer struct {
	cfg                      *config.Config
	log                      *zap.Logger
	db                       *pgxpool.Pool
	kafka                    *kafka.Kafka
	messageConsumer          *events.MessageConsumerEvent
	validator                *validator.Validate
	orderHandler             *handlers.OrdersHandler
	statusHandler            *handlers.GetStatusHandler
	ordersService            servInterfaces.OrderService
	statusService            servInterfaces.StatusService
	messageTrackerRepository repoInterfaces.MessageTrackerRepository
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

func (di *DiContainer) Kafka(_ context.Context) *kafka.Kafka {
	if di.kafka == nil {
		client, err := kafka.NewKafkaClient(di.cfg.Kafka, di.log)
		if err != nil {
			os.Exit(1)
		}

		di.kafka = client
	}

	return di.kafka
}

func (di *DiContainer) MessageConsumer(ctx context.Context) *events.MessageConsumerEvent {
	if di.messageConsumer == nil {
		di.messageConsumer = events.NewMessageConsumerEvent(di.log, di.Kafka(ctx), di.OrdersService(ctx), di.Validator(ctx))
	}

	return di.messageConsumer
}

func (di *DiContainer) OrdersService(ctx context.Context) servInterfaces.OrderService {
	if di.ordersService == nil {
		di.ordersService = order_service.NewOrderService(di.MessageTrackerRepository(ctx), di.Kafka(ctx))
	}

	return di.ordersService
}

func (di *DiContainer) StatusService(ctx context.Context) servInterfaces.StatusService {
	if di.statusService == nil {
		di.statusService = statuses_service.NewStatusService(di.MessageTrackerRepository(ctx))
	}

	return di.statusService
}

func (di *DiContainer) MessageTrackerRepository(ctx context.Context) repoInterfaces.MessageTrackerRepository {
	if di.messageTrackerRepository == nil {
		di.messageTrackerRepository = mesTrackerRepo.NewMessageTrackerRepository(di.DB(ctx))
	}

	return di.messageTrackerRepository
}

func (di *DiContainer) Validator(_ context.Context) *validator.Validate {
	if di.validator == nil {
		di.validator = validator.New()
	}

	return di.validator
}

func (di *DiContainer) OrdersHandler(ctx context.Context) *handlers.OrdersHandler {
	if di.orderHandler == nil {
		di.orderHandler = handlers.NewOrdersHandler(di.log, di.OrdersService(ctx), di.Validator(ctx))
	}

	return di.orderHandler
}

func (di *DiContainer) GetStatusHandler(ctx context.Context) *handlers.GetStatusHandler {
	if di.statusHandler == nil {
		di.statusHandler = handlers.NewGetStatusHandler(di.log, di.StatusService(ctx))
	}

	return di.statusHandler
}
