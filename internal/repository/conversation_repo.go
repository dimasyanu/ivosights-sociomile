package repository

import (
	"github.com/dimasyanu/ivosights-sociomile/domain"
	"github.com/google/uuid"
)

type ConversationRepository interface {
	GetByID(id uuid.UUID) (*domain.ConversationEntity, error)
	GetByTenantAndCustomer(tenantID uint, customerID uuid.UUID) (*domain.ConversationEntity, error)
	Create(conversation *domain.ConversationEntity) (uuid.UUID, error)
	UpdateStatus(id uuid.UUID, status string) error
	Delete(id uuid.UUID) error
}
