package service

import (
	"fmt"
	"time"

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
		if err != repository.ErrNotFound {
			return nil, err
		}
		// If not found, create a new conversation
		return s.Create(tID, custID)
	}
	return convEntity.ToDto(), nil
}

func (s *ConversationService) Create(tID uint, custID uuid.UUID) (*domain.Conversation, error) {
	convEntity := &domain.ConversationEntity{
		ID:         uuid.New(),
		TenantID:   tID,
		CustomerID: custID,
		Status:     "open",
		CreatedAt:  time.Now(),
	}
	_, err := s.convRepo.Create(convEntity)
	if err != nil {
		return nil, err
	}

	// Send to message queue for async processing
	go func() {
		err = s.mq.PublishMessage("conversation_created", []byte(convEntity.ID.String()))
		if err != nil {
			fmt.Printf("Failed to publish message for conversation creation: %v\n", err)
			return
		}
	}()

	return convEntity.ToDto(), nil
}

func (s *ConversationService) UpdateStatus(id uuid.UUID, status string) error {
	err := s.convRepo.UpdateStatus(id, status)
	if err != nil {
		return err
	}

	// Send to message queue for async processing
	go func() {
		err = s.mq.PublishMessage("conversation_status_updated", []byte(id.String()))
		if err != nil {
			fmt.Printf("Failed to publish message for conversation status update: %v\n", err)
			return
		}
	}()

	return nil
}
