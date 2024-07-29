package message_tracker_repository

import _ "embed"

var (
	//go:embed sql/save_message.sql
	SaveMessage string

	//go:embed sql/update_message_status_by_request_id.sql
	UpdateMessageStatusByRequestId string
)
