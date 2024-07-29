UPDATE message_tracker SET status = $1, updated_at = NOW() WHERE request_id = $2;
