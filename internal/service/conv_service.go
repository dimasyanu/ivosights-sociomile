package service

import (
	"encoding/json"
	"fmt"
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

func (s *ConversationService) GetByID(id uuid.UUID) (*domain.Conversation, error) {
	convEntity, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return convEntity.ToDto(), nil
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

	// Send to message queue for async processing
	go func() {
		data := &infra.ConversationStatusUpdatedMessage{
			ConvID: id,
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
