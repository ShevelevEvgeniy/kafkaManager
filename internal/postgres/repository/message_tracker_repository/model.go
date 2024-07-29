package message_tracker_repository

import "time"

type Model struct {
	ID        int64     `json:"id"`
	RequestId string    `json:"request_id"`
	Message   []byte    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
