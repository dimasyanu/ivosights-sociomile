package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dimasyanu/ivosights-sociomile/internal/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/domain/constant"
	"github.com/dimasyanu/ivosights-sociomile/internal/domain/repo"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra"
	"github.com/google/uuid"
)

type TicketService struct {
	repo    repo.TicketRepository
	convSvc *ConversationService
	mq      infra.QueueClient
}

func NewTicketService(ticketRepo repo.TicketRepository, convSvc *ConversationService, mq infra.QueueClient) *TicketService {
	return &TicketService{repo: ticketRepo, convSvc: convSvc, mq: mq}
}

func (s *TicketService) GetList(f *domain.TicketFilter) (*domain.Paginated[domain.Ticket], error) {
	ticketEntities, total, err := s.repo.GetList(f)
	if err != nil {
		return nil, err
	}
	tickets := []domain.Ticket{}
	for _, entity := range ticketEntities {
		tickets = append(tickets, *entity.ToDto())
	}
	return &domain.Paginated[domain.Ticket]{
		Items:    tickets,
		Total:    total,
		Page:     f.Page,
		PageSize: f.PageSize,
	}, nil
}

func (s *TicketService) GetByID(id uuid.UUID) (*domain.Ticket, error) {
	ticketEntity, err := s.repo.GetByID(id)
	if err != nil {
		idStr := id.String()
		fmt.Printf("%s", idStr)
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

func (s *TicketService) ticketCreatedEvent(t *domain.Ticket) {
	// Send to message queue for async processing
	data := &infra.TicketCreatedMessage{
		ConvID:   t.ConversationID,
		TicketID: t.ID,
		Status:   t.Status,
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Failed to marshal ticket created message: %v\n", err)
		return
	}
	err = s.mq.Publish("ticket_created", bytes)
	if err != nil {
		fmt.Printf("Failed to publish message for ticket created: %v\n", err)
		return
	}
}

func (s *TicketService) Create(convID uuid.UUID, title, description string, priority int8, pic string) (*domain.Ticket, error) {
	// Check if conversation exists
	conv, err := s.convSvc.GetByID(convID)
	if err != nil {
		return nil, err
	}

	// Check if ticket already exists for this conversation
	t, err := s.GetByConversationID(convID)
	if err != nil && !errors.Is(err, repo.ErrNotFound) {
		return nil, err
	}
	if t != nil {
		return nil, errors.New("ticket already exists for this conversation")
	}

	e := &domain.TicketEntity{
		ConversationID: convID,
		TenantID:       conv.TenantID,
		Title:          title,
		Description:    description,
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		CreatedBy:      pic,
		UpdatedAt:      time.Now(),
		UpdatedBy:      pic,
		Priority:       priority,
		Status:         constant.TicketStatusOpen,
	}

	createdEntity, err := s.repo.Create(e)
	if err != nil {
		return nil, err
	}

	// Send event to message queue for async processing
	go s.ticketCreatedEvent(createdEntity.ToDto())

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
