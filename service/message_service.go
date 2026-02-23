package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dimasyanu/ivosights-sociomile/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra"
	"github.com/dimasyanu/ivosights-sociomile/internal/repository"
	"github.com/google/uuid"
)

type MessageService struct {
	mq      infra.QueueClient
	convSvc *ConversationService
	repo    repository.MessageRepository
}

func NewMessageService(convSvc *ConversationService, repo repository.MessageRepository, mq infra.QueueClient) *MessageService {
	return &MessageService{
		convSvc: convSvc,
		repo:    repo,
		mq:      mq,
	}
}

func (s *MessageService) GetMessages(convID uuid.UUID, offset int64, limit int) ([]*domain.Message, int64, error) {
	_, err := s.convSvc.GetByID(convID)
	if err != nil {
		return nil, 0, err
	}

	messages, total, err := s.repo.GetByConversationID(convID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	// Map MessageEntity to Message
	messageDtos := make([]*domain.Message, len(messages))
	for i, msg := range messages {
		messageDtos[i] = msg.ToDto()
	}

	return messageDtos, total, nil
}

func (s *MessageService) CreateMessage(tID uint, custID uuid.UUID, senderType string, message string) (*domain.Message, error) {
	// Check if conversation exists for the given tenant and customer
	conv, err := s.convSvc.GetByTenantAndCustomer(tID, custID)
	if err != nil {
		if err != repository.ErrNotFound {
			return nil, err
		}
		// If not found, create a new conversation
		conv, err = s.createNewConversation(tID, custID, message)
		if err != nil {
			return nil, err
		}
	}

	messageEntity := &domain.MessageEntity{
		ConversationID: conv.ID,
		SenderType:     senderType,
		Message:        message,
		CreatedAt:      time.Now(),
	}

	messageEntity.ConversationID = conv.ID
	id, err := s.repo.Create(messageEntity)
	if err != nil {
		return nil, err
	}

	messageEntity.ID = id
	return messageEntity.ToDto(), nil
}

func (s *MessageService) createNewConversation(tID uint, custID uuid.UUID, msg string) (*domain.Conversation, error) {
	// If not found, create a new conversation
	conv, err := s.convSvc.Create(tID, custID)
	if err != nil {
		return nil, err
	}

	// Send to message queue for async processing
	go func() {
		data := &infra.ConversationCreatedMessage{
			TenantID: tID,
			CustID:   custID,
			ConvID:   conv.ID,
			Message:  msg,
		}
		bytes, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("Failed to marshal conversation created message: %v\n", err)
			return
		}
		err = s.mq.Publish("conversation_created", bytes)
		if err != nil {
			fmt.Printf("Failed to publish message for conversation creation: %v\n", err)
			return
		}
	}()

	return conv, nil
}

func (s *MessageService) DeleteMessage(id uuid.UUID) error {
	return s.repo.Delete(id)
}
