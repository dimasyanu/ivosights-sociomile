package mysqlrepo

import (
	"database/sql"

	"github.com/dimasyanu/ivosights-sociomile/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/repository"
	"github.com/google/uuid"
)

type ConversationMySqlRepository struct {
	db *sql.DB
}

// Create implements [repository.ConversationRepository].
func (c *ConversationMySqlRepository) Create(conversation *domain.ConversationEntity) (uuid.UUID, error) {
	panic("unimplemented")
}

// Delete implements [repository.ConversationRepository].
func (c *ConversationMySqlRepository) Delete(id uuid.UUID) error {
	panic("unimplemented")
}

// GetByID implements [repository.ConversationRepository].
func (c *ConversationMySqlRepository) GetByID(id uuid.UUID) (*domain.ConversationEntity, error) {
	panic("unimplemented")
}

// GetByTenantAndCustomer implements [repository.ConversationRepository].
func (c *ConversationMySqlRepository) GetByTenantAndCustomer(tenantID uint, customerID uuid.UUID) (*domain.ConversationEntity, error) {
	panic("unimplemented")
}

// UpdateStatus implements [repository.ConversationRepository].
func (c *ConversationMySqlRepository) UpdateStatus(id uuid.UUID, status string) error {
	panic("unimplemented")
}

func NewConversationRepository(db *sql.DB) repository.ConversationRepository {
	// Initialize your database connection here and assign it to the db field
	return &ConversationMySqlRepository{
		db: db, // Replace with actual DB connection
	}
}
