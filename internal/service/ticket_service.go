package service

import (
	"database/sql"
	"errors"

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

func (s *TicketService) Create(entity *domain.TicketEntity) (*domain.Ticket, error) {
	createdEntity, err := s.repo.Create(entity)
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
