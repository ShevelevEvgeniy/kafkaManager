package kafka

import (
	"context"

	"github.com/ShevelevEvgeniy/kafkaManager/config"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Kafka struct {
	producer    *kafka.Producer
	consumer    *kafka.Consumer
	adminClient *kafka.AdminClient
	cfg         config.Kafka
	log         *zap.Logger
}

func NewKafkaClient(cfg config.Kafka, log *zap.Logger) (*Kafka, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":  cfg.Broker,
		"acks":               "all",
		"retries":            cfg.Retries,
		"request.timeout.ms": cfg.Timeout,
	})
	if err != nil {
		log.Error("failed to create kafka producer", zap.Error(err))
		return nil, errors.Wrap(err, "failed to create kafka producer")
	}

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{"bootstrap.servers": cfg.Broker, "group.id": "messaggio", "auto.offset.reset": "earliest"})
	if err != nil {
		log.Error("failed to create kafka consumer", zap.Error(err))
		return nil, errors.Wrap(err, "failed to create kafka consumer")
	}

	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": cfg.Broker})
	if err != nil {
		log.Error("failed to create kafka admin client", zap.Error(err))
		return nil, errors.Wrap(err, "failed to create kafka admin client")
	}

	log.Info("created kafka producer", zap.String("broker", cfg.Broker), zap.String("topic", cfg.Topic))

	client := &Kafka{
		producer:    producer,
		consumer:    consumer,
		adminClient: adminClient,
		cfg:         cfg,
		log:         log,
	}

	go client.handleDeliveryReports()

	return client, nil
}

func (c *Kafka) SendMessage(ctx context.Context, key, value []byte, topic string) error {
	deliveryChan := make(chan kafka.Event, 1)
	defer close(deliveryChan)

	err := c.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          value,
	}, deliveryChan)
	if err != nil {
		return errors.Wrap(err, "failed to send message to kafka")
	}

	select {
	case event := <-deliveryChan:
		if msg, ok := event.(*kafka.Message); ok {
			if msg.TopicPartition.Error != nil {
				return errors.Wrap(msg.TopicPartition.Error, "failed to deliver message")
			}

			return nil
		}

		return errors.Wrap(err, "unexpected event type received")
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Kafka) ListenToTopic(ctx context.Context, topic string) (<-chan *kafka.Message, <-chan error) {
	messageChan := make(chan *kafka.Message)
	errChan := make(chan error, 1)

	go func() {
		defer close(messageChan)
		defer close(errChan)

		for {
			select {
			case <-ctx.Done():
				c.log.Info("shutting down message consumer")
				errChan <- ctx.Err()
				return
			default:
				msg, err := c.consumer.ReadMessage(-1)
				if err != nil {
					var kafkaError kafka.Error
					if errors.As(err, &kafkaError) && kafkaError.Code() == kafka.ErrPartitionEOF {
						continue
					}

					c.log.Error("failed to read message from topic", zap.String("topic", topic), zap.Error(err))
					errChan <- errors.Wrap(err, "failed to read message from topic")
					return
				}

				messageChan <- msg
			}
		}
	}()

	return messageChan, nil
}

func (c *Kafka) SubscribeToTopics(ctx context.Context, topic string) error {
	err := c.CreateTopic(ctx, []string{topic})
	if err != nil {
		c.log.Error("failed to create topic", zap.String("topic", topic), zap.Error(err))
		return errors.Wrap(err, "failed to create topic")
	}

	err = c.consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		c.log.Error("failed to subscribe to topic", zap.String("topic", topic), zap.Error(err))
		return errors.Wrap(err, "failed to subscribe to topic")
	}

	return nil
}

func (c *Kafka) handleDeliveryReports() {
	for e := range c.producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				c.log.Error("Delivery failed", zap.Error(ev.TopicPartition.Error))
			} else {
				c.log.Error("Delivered message", zap.String("value", string(ev.Value)))
			}
		}
	}
}

func (c *Kafka) CreateTopic(ctx context.Context, topics []string) error {
	for _, topic := range topics {
		exists, err := c.topicExists(ctx, topic)
		if err != nil {
			return errors.Wrap(err, "failed to check if topic exists")
		}

		if exists {
			return nil
		}

		topicSpec := kafka.TopicSpecification{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}

		_, err = c.adminClient.CreateTopics(ctx, []kafka.TopicSpecification{topicSpec})
		if err != nil {
			return errors.Wrap(err, "failed to create topic")
		}
	}

	return nil
}

func (c *Kafka) topicExists(_ context.Context, topic string) (bool, error) {
	metadata, err := c.adminClient.GetMetadata(&topic, false, 500)
	if err != nil {
		return false, errors.Wrap(err, "failed to get metadata")
	}

	_, exists := metadata.Topics[topic]
	return exists, nil
}

func (c *Kafka) Close() {
	_ = c.consumer.Close()
	c.producer.Close()
	c.adminClient.Close()
}
