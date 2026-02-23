package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dimasyanu/ivosights-sociomile/constant"
	"github.com/dimasyanu/ivosights-sociomile/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra"
	"github.com/dimasyanu/ivosights-sociomile/internal/repository"
	"github.com/google/uuid"
)

type ConversationService struct {
	convRepo repository.ConversationRepository
	userSvc  *UserService
	mq       infra.QueueClient
}

func NewConversationService(convRepo repository.ConversationRepository, userSvc *UserService, mq infra.QueueClient) *ConversationService {
	return &ConversationService{convRepo: convRepo, userSvc: userSvc, mq: mq}
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

func (s *ConversationService) Create(tID uint, custID uuid.UUID) (*domain.Conversation, error) {
	convEntity := &domain.ConversationEntity{
		ID:         uuid.New(),
		TenantID:   tID,
		CustomerID: custID,
		Status:     constant.ConvStatusOpen,
		CreatedAt:  time.Now(),
	}
	_, err := s.convRepo.Create(convEntity)
	if err != nil {
		return nil, err
	}

	return convEntity.ToDto(), nil
}

func (s *ConversationService) AssignConversation(id uuid.UUID) error {
	conv, err := s.GetByID(id)
	if err != nil {
		return err
	}

	// If conversation is already assigned, do nothing
	if conv.AssignedAgentID != nil {
		return nil
	}

	err := s.userSvc.GetAvailableAgent()

	return nil
}

func (s *ConversationService) UpdateStatus(id uuid.UUID, status string) error {
	err := s.convRepo.UpdateStatus(id, status)
	if err != nil {
		return err
	}

	// Send to message queue for async processing
	go func() {
		data := &infra.ConversationStatusUpdatedMessage{
			ConvID: id.String(),
			Status: status,
		}
		bytes, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("Failed to marshal conversation status updated message: %v\n", err)
			return
		}
		err = s.mq.Publish("conversation_status_updated", bytes)
		if err != nil {
			fmt.Printf("Failed to publish message for conversation status update: %v\n", err)
			return
		}
	}()

	return nil
}
