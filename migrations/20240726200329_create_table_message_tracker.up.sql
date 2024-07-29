CREATE TYPE message_status AS ENUM ('sent', 'failed', 'succeeded');

CREATE TABLE message_tracker (
    request_id VARCHAR(255) PRIMARY KEY NOT NULL,
    message BYTEA NOT NULL,
    status message_status NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_message_tracker_status ON message_tracker(status);