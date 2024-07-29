package convertor

import (
	"github.com/ShevelevEvgeniy/kafkaManager/internal/dto"
	msgTracRepo "github.com/ShevelevEvgeniy/kafkaManager/internal/postgres/repository/message_tracker_repository"
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
