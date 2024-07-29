package message_tracker_repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewMessageTrackerRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) SaveMessage(ctx context.Context, model Model) error {
	query := SaveMessage

	_, err := r.db.Exec(ctx, query, model.RequestId, model.Message, model.Status)
	if err != nil {
		return errors.Wrap(err, "failed to save message")
	}

	return nil
}

func (r *Repository) UpdateMessageStatusByRequestId(ctx context.Context, model Model) error {
	query := UpdateMessageStatusByRequestId

	_, err := r.db.Exec(ctx, query, model.Status, model.RequestId)
	if err != nil {
		return errors.Wrap(err, "failed to update message status")
	}
	return nil
}
