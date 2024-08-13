package convertor

import (
	"github.com/ShevelevEvgeniy/kafkaManager/internal/dto"
	msgTracRepo "github.com/ShevelevEvgeniy/kafkaManager/internal/postgres/repository/message_tracker_repository"
	"github.com/ShevelevEvgeniy/kafkaManager/pkg/event_dispatcher"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func OrderDtoToTrackingModel(order []byte, requestId string, status string) msgTracRepo.Model {
	return msgTracRepo.Model{
		RequestId: requestId,
		Message:   order,
		Status:    status,
	}
}

func OrderMessageDtoToTrackingModel(order dto.OrderMessageResponse) msgTracRepo.Model {
	return msgTracRepo.Model{
		RequestId: order.RequestId,
		Status:    order.Status,
	}
}

func KafkaMessageToMessage(msg *kafka.Message) event_dispatcher.Message {
	return event_dispatcher.Message{
		Key:   msg.Key,
		Value: msg.Value,
	}
}
