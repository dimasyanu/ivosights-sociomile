-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS conversations (
    id BINARY(16) NOT NULL PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    tenant_id INT NOT NULL,
    customer_id BINARY(16) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'open',
    assigned_agent_id BINARY(16) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_customer_id (customer_id),
    INDEX idx_assigned_agent_id (assigned_agent_id),

    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (assigned_agent_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS conversations;
-- +goose StatementEnd
