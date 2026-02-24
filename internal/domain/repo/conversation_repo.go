package repo

import (
	"github.com/dimasyanu/ivosights-sociomile/internal/domain"
	"github.com/google/uuid"
)

type ConversationRepository interface {
	GetList(filter *domain.ConversationFilter) ([]domain.ConversationEntity, uint64, error)
	GetByID(id uuid.UUID) (*domain.ConversationEntity, error)
	GetByTenantAndCustomer(tenantID uint, customerID uuid.UUID) (*domain.ConversationEntity, error)
	Create(conversation *domain.ConversationEntity) (uuid.UUID, error)
	UpdateStatus(id uuid.UUID, status string) error
	UpdateAssignment(conv *domain.ConversationEntity, agentID uuid.UUID) error
}
