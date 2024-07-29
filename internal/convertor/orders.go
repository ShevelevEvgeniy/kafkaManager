package convertor

import (
	msgTracRepo "github.com/ShevelevEvgeniy/kafkaManager/internal/postgres/repository/message_tracker_repository"
)

func OrderDtoToTrackingModel(order []byte, requestId string, status string) msgTracRepo.Model {
	return msgTracRepo.Model{
		RequestId: requestId,
		Message:   order,
		Status:    status,
	}
}
