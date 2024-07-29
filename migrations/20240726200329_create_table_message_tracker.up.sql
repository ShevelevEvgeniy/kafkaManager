CREATE TYPE message_status AS ENUM ('sent', 'failed', 'succeeded');

CREATE TABLE message_tracker (
    request_id VARCHAR(255) NOT NULL AS IDENTITY PRIMARY KEY,
    message BYTEA NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_message_tracker_request_id ON message_tracker(request_id);
CREATE INDEX idx_message_tracker_status ON message_tracker(status);