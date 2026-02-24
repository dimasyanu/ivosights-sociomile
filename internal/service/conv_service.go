package service

import (
	"time"

	"github.com/dimasyanu/ivosights-sociomile/internal/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/domain/constant"
	"github.com/dimasyanu/ivosights-sociomile/internal/domain/repo"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra"
	"github.com/google/uuid"
)

type ConversationService struct {
	repo repo.ConversationRepository
	mq   infra.QueueClient
}

func NewConversationService(convRepo repo.ConversationRepository, mq infra.QueueClient) *ConversationService {
	return &ConversationService{repo: convRepo, mq: mq}
}

func (s *ConversationService) GetList(filter *domain.ConversationFilter) (*domain.Paginated[domain.Conversation], error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 25
	}
	convEntities, total, err := s.repo.GetList(filter)
	if err != nil {
		return nil, err
	}

	convs := make([]domain.Conversation, len(convEntities))
	for i, e := range convEntities {
		convs[i] = *e.ToDto()
	}

	return &domain.Paginated[domain.Conversation]{
		Items:    convs,
		Total:    total,
		Page:     filter.Page,
		PageSize: filter.PageSize,
	}, nil
}

func (s *ConversationService) GetByID(id uuid.UUID) (*domain.ConversationDetail, error) {
	convEntity, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return convEntity.ToDetailDto(), nil
}

func (s *ConversationService) GetByTenantAndCustomer(tID uint, custID uuid.UUID) (*domain.Conversation, error) {
	convEntity, err := s.repo.GetByTenantAndCustomer(tID, custID)
	if err != nil {
		return nil, err
	}
	return convEntity.ToDto(), nil
}

func (s *ConversationService) Create(tID uint, custID uuid.UUID) (*domain.Conversation, error) {
	convEntity := &domain.ConversationEntity{
		ID:         uuid.New(),
		TenantID:   tID,
		CustomerID: custID,
		Status:     constant.ConvStatusOpen,
		CreatedAt:  time.Now(),
	}
	_, err := s.repo.Create(convEntity)
	if err != nil {
		return nil, err
	}

	return convEntity.ToDto(), nil
}

func (s *ConversationService) AssignConversation(id uuid.UUID, agentID uuid.UUID) error {
	conv, err := s.GetByID(id)
	if err != nil {
		return err
	}

	// If conversation is already assigned, do nothing
	if conv.AssignedAgentID != nil {
		return nil
	}

	data := &domain.ConversationEntity{
		ID:       conv.ID,
		TenantID: conv.TenantID,
	}
	return s.repo.UpdateAssignment(data, agentID)
}

func (s *ConversationService) UpdateStatus(id uuid.UUID, status string) error {
	err := s.repo.UpdateStatus(id, status)
	if err != nil {
		return err
	}

	return nil
}

func (s *ConversationService) Delete(id uuid.UUID) error {
	return s.repo.UpdateStatus(id, constant.ConvStatusClosed)
}
