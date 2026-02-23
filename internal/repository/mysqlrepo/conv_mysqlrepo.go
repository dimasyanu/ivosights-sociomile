package mysqlrepo

import (
	"context"
	"database/sql"

	"github.com/dimasyanu/ivosights-sociomile/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/repository"
	"github.com/google/uuid"
)

type ConversationMySqlRepository struct {
	db *sql.DB
}

// Create implements [repository.ConversationRepository].
func (r *ConversationMySqlRepository) Create(c *domain.ConversationEntity) (uuid.UUID, error) {
	id := uuid.New()
	pairs := map[string]any{
		"id":          id,
		"tenant_id":   c.TenantID,
		"customer_id": c.CustomerID,
		"status":      c.Status,
		"created_at":  c.CreatedAt,
	}
	cols, slots, vals := MapForCreate(pairs)

	query := "INSERT INTO conversations (" + cols + ") VALUES (" + slots + ")"
	_, err := r.db.ExecContext(context.Background(), query, vals...)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

// GetByID implements [repository.ConversationRepository].
func (r *ConversationMySqlRepository) GetByID(id uuid.UUID) (*domain.ConversationEntity, error) {
	panic("unimplemented")
}

// GetByTenantAndCustomer implements [repository.ConversationRepository].
func (r *ConversationMySqlRepository) GetByTenantAndCustomer(tenantID uint, customerID uuid.UUID) (*domain.ConversationEntity, error) {
	var conversation domain.ConversationEntity
	err := r.db.QueryRow("SELECT id, tenant_id, customer_id, status FROM conversations WHERE tenant_id = ? AND customer_id = ?", tenantID, customerID).Scan(&conversation.ID, &conversation.TenantID, &conversation.CustomerID, &conversation.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &conversation, nil
}

// UpdateStatus implements [repository.ConversationRepository].
func (r *ConversationMySqlRepository) UpdateStatus(id uuid.UUID, status string) error {
	panic("unimplemented")
}

func NewConversationRepository(db *sql.DB) repository.ConversationRepository {
	// Initialize your database connection here and assign it to the db field
	return &ConversationMySqlRepository{
		db: db, // Replace with actual DB connection
	}
}
