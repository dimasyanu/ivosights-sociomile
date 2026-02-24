-- +goose Up
-- +goose StatementBegin
CREATE TABLE messages (
    id BINARY(16) NOT NULL PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    conversation_id BINARY(16) NOT NULL,
    sender_type VARCHAR(20) NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,

    INDEX idx_conversation_id (conversation_id),
    INDEX idx_sender_type (sender_type),

    FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS messages;
-- +goose StatementEnd
