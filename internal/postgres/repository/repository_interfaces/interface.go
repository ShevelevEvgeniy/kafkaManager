package repository_interfaces

import (
	"context"

	mesTrackModel "github.com/ShevelevEvgeniy/kafkaManager/internal/postgres/repository/message_tracker_repository"
)

type MessageTrackerRepository interface {
	SaveMessage(ctx context.Context, model mesTrackModel.Model) error
	UpdateMessageStatusByRequestId(ctx context.Context, model mesTrackModel.Model) error
}
