package service

import (
	"github.com/dimasyanu/ivosights-sociomile/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra"
	"github.com/dimasyanu/ivosights-sociomile/internal/repository"
	"github.com/google/uuid"
)

type ConversationService struct {
	convRepo repository.ConversationRepository
	mq       infra.QueueClient
}

func NewConversationService(convRepo repository.ConversationRepository, mq infra.QueueClient) *ConversationService {
	return &ConversationService{convRepo: convRepo, mq: mq}
}

func (s *ConversationService) GetByID(id uuid.UUID) (*domain.Conversation, error) {
	convEntity, err := s.convRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return convEntity.ToDto(), nil
}

func (s *ConversationService) GetByTenantAndCustomer(tID uint, custID uuid.UUID) (*domain.Conversation, error) {
	convEntity, err := s.convRepo.GetByTenantAndCustomer(tID, custID)
	if err != nil {
		return nil, err
	}
	return convEntity.ToDto(), nil
}

func (s *ConversationService) Create(tID uint, custID uuid.UUID) (uuid.UUID, error) {
	convEntity := &domain.ConversationEntity{
		ID:         uuid.New(),
		TenantID:   tID,
		CustomerID: custID,
		Status:     "open",
	}
	id, err := s.convRepo.Create(convEntity)
	if err != nil {
		return uuid.Nil, err
	}

	// Send to message queue for async processing
	err = s.mq.PublishMessage("conversation_created", []byte(convEntity.ID.String()))
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (s *ConversationService) UpdateStatus(id uuid.UUID, status string) error {
	err := s.convRepo.UpdateStatus(id, status)
	if err != nil {
		return err
	}

	// Send to message queue for async processing
	err = s.mq.PublishMessage("conversation_status_updated", []byte(id.String()))
	if err != nil {
		return err
	}

	return nil
}
