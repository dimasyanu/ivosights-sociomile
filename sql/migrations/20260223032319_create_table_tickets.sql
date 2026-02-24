-- +goose Up
-- +goose StatementBegin
CREATE TABLE tickets (
    id BINARY(16) NOT NULL PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID(), 1)),
    tenant_id INT NOT NULL,
    conversation_id BINARY(16) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL,
    priority TINYINT NOT NULL,
    assigned_agent_id BINARY(16) NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(255) NOT NULL,
    deleted_at TIMESTAMP NULL,
    deleted_by VARCHAR(255) NULL,

    FOREIGN KEY (assigned_agent_id) REFERENCES users(id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),

    UNIQUE KEY unique_conversation (conversation_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_conversation_id (conversation_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tickets;
-- +goose StatementEnd
