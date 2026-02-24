package repo

import (
	"github.com/dimasyanu/ivosights-sociomile/internal/domain"
	"github.com/google/uuid"
)

type TicketRepository interface {
	GetList(f *domain.TicketFilter) ([]*domain.TicketEntity, uint64, error)
	GetByID(id uuid.UUID) (*domain.TicketEntity, error)
	GetByConversationID(convID uuid.UUID) (*domain.TicketEntity, error)
	Create(e *domain.TicketEntity) (*domain.TicketEntity, error)
	Update(e *domain.TicketEntity) (*domain.TicketEntity, error)
	UpdateStatus(id uuid.UUID, status string) (*domain.TicketEntity, error)
}
