package statuses_service

import (
	"context"

	repoInterfaces "github.com/ShevelevEvgeniy/kafkaManager/internal/postgres/repository/repository_interfaces"
	"github.com/pkg/errors"
)

type StatusService struct {
	repository repoInterfaces.MessageTrackerRepository
}

func NewStatusService(repository repoInterfaces.MessageTrackerRepository) *StatusService {
	return &StatusService{
		repository: repository,
	}
}

func (s *StatusService) GetStatus(ctx context.Context, requestId string) (string, error) {
	status, err := s.repository.GetMessageByRequestId(ctx, requestId)
	if err != nil {
		return "", errors.Wrap(err, "failed to get message by request id")
	}

	return status, nil
}
