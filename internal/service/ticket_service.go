package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/dimasyanu/ivosights-sociomile/internal/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/domain/repo"
	"github.com/google/uuid"
)

type TicketService struct {
	repo repo.TicketRepository
}

func NewTicketService(ticketRepo repo.TicketRepository) *TicketService {
	return &TicketService{repo: ticketRepo}
}

func (s *TicketService) GetList(f *domain.TicketFilter) ([]*domain.Ticket, uint64, error) {
	ticketEntities, total, err := s.repo.GetList(f)
	if err != nil {
		return nil, 0, err
	}
	tickets := make([]*domain.Ticket, len(ticketEntities))
	for i, entity := range ticketEntities {
		tickets[i] = entity.ToDto()
	}
	return tickets, total, nil
}

func (s *TicketService) GetByID(id uuid.UUID) (*domain.Ticket, error) {
	ticketEntity, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repo.ErrNotFound
		}
		return nil, err
	}
	if ticketEntity == nil {
		return nil, repo.ErrNotFound
	}
	return ticketEntity.ToDto(), nil
}

func (s *TicketService) GetByConversationID(convID uuid.UUID) (*domain.Ticket, error) {
	ticketEntitiy, err := s.repo.GetByConversationID(convID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repo.ErrNotFound
		}
		return nil, err
	}
	if ticketEntitiy == nil {
		return nil, repo.ErrNotFound
	}
	return ticketEntitiy.ToDto(), nil
}

func (s *TicketService) Create(e *domain.TicketEntity, pic string) (*domain.Ticket, error) {
	t, err := s.GetByConversationID(e.ConversationID)
	if err != nil && !errors.Is(err, repo.ErrNotFound) {
		return nil, err
	}
	if t != nil {
		return nil, errors.New("ticket already exists for this conversation")
	}

	e.ID = uuid.New()
	e.CreatedAt = time.Now()
	e.CreatedBy = pic
	e.UpdatedAt = time.Now()
	e.UpdatedBy = pic

	createdEntity, err := s.repo.Create(e)
	if err != nil {
		return nil, err
	}
	return createdEntity.ToDto(), nil
}

func (s *TicketService) Update(entity *domain.TicketEntity) (*domain.Ticket, error) {
	updatedEntity, err := s.repo.Update(entity)
	if err != nil {
		return nil, err
	}
	return updatedEntity.ToDto(), nil
}

func (s *TicketService) UpdateStatus(id uuid.UUID, status string) (*domain.Ticket, error) {
	ticket, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		return nil, repo.ErrNotFound
	}
	updatedEntity, err := s.repo.UpdateStatus(ticket.ID, status)
	if err != nil {
		return nil, err
	}
	return updatedEntity.ToDto(), nil
}

func (s *TicketService) Delete(id uuid.UUID) error {
	ticket, err := s.GetByID(id)
	if err != nil {
		return err
	}
	now := time.Now()
	ticket.DeletedAt = &now
	_, err = s.repo.Update(&domain.TicketEntity{
		ID:        ticket.ID,
		DeletedAt: ticket.DeletedAt,
	})
	if err != nil {
		return err
	}
	return nil
}
