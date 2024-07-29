package config

import "time"

type Kafka struct {
	Broker  string        `envconfig:"KAFKA_BROKER" env-required:"true"`
	Topic   string        `envconfig:"KAFKA_TOPIC" env-required:"true"`
	Retries int           `envconfig:"KAFKA_RETRIES" default:"3"`
	Timeout time.Duration `envconfig:"KAFKA_TIMEOUT" default:"10s"`
}
