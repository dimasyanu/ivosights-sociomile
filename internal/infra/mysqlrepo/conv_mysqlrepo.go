package mysqlrepo

import (
	"context"
	"database/sql"
	"log"

	"github.com/dimasyanu/ivosights-sociomile/internal/domain"
	repository "github.com/dimasyanu/ivosights-sociomile/internal/domain/repo"
	"github.com/google/uuid"
)

type ConversationMySqlRepository struct {
	db *sql.DB
}

func NewConversationRepository(db *sql.DB) repository.ConversationRepository {
	// Initialize your database connection here and assign it to the db field
	return &ConversationMySqlRepository{
		db: db, // Replace with actual DB connection
	}
}

// Create implements [repository.ConversationRepository].
func (r *ConversationMySqlRepository) Create(c *domain.ConversationEntity) (uuid.UUID, error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	pairs := map[string]any{
		"id":          c.ID,
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

	return c.ID, nil
}

// GetByID implements [repository.ConversationRepository].
func (r *ConversationMySqlRepository) GetByID(id uuid.UUID) (*domain.ConversationEntity, error) {
	conversation := &domain.ConversationEntity{}
	query := "SELECT id, tenant_id, customer_id, status FROM conversations WHERE id = UUID_TO_BIN(?)"
	err := r.db.QueryRow(query, id).Scan(&conversation.ID, &conversation.TenantID, &conversation.CustomerID, &conversation.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return conversation, nil
}

// GetByTenantAndCustomer implements [repository.ConversationRepository].
func (r *ConversationMySqlRepository) GetByTenantAndCustomer(tenantID uint, customerID uuid.UUID) (*domain.ConversationEntity, error) {
	var conversation domain.ConversationEntity
	err := r.db.QueryRow("SELECT id, tenant_id, customer_id, assigned_agent_id, status FROM conversations WHERE tenant_id = ? AND customer_id = UUID_TO_BIN(?)", tenantID, customerID.String()).Scan(&conversation.ID, &conversation.TenantID, &conversation.CustomerID, &conversation.AssignedAgentID, &conversation.Status)
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
	query := "UPDATE conversations SET status = ? WHERE id = UUID_TO_BIN(?)"
	_, err := r.db.ExecContext(context.Background(), query, status, id)
	return err
}

// UpdateAssignment implements [repository.ConversationRepository].
func (r *ConversationMySqlRepository) UpdateAssignment(conv *domain.ConversationEntity, agentID uuid.UUID) error {
	query := "UPDATE conversations SET assigned_agent_id = UUID_TO_BIN(?) WHERE id = UUID_TO_BIN(?)"
	_, err := r.db.ExecContext(context.Background(), query, agentID, conv.ID)
	log.Printf("AgentID: %s\n", agentID.String())
	return err
}
