package message_tracker_repository

import _ "embed"

var (
	//go:embed sql/save_message.sql
	SaveMessage string
)
