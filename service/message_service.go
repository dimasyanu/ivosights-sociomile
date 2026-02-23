package service

import (
	"time"

	"github.com/dimasyanu/ivosights-sociomile/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra"
	"github.com/dimasyanu/ivosights-sociomile/internal/repository"
	"github.com/google/uuid"
)

type MessageService struct {
	convSvc *ConversationService
	repo    repository.MessageRepository
}

func NewMessageService(convSvc *ConversationService, repo repository.MessageRepository, mq infra.QueueClient) *MessageService {
	return &MessageService{
		convSvc: convSvc,
		repo:    repo,
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

func (s *MessageService) CreateMessage(tID uint, custID uuid.UUID, senderType string, message string) (uuid.UUID, error) {
	// Check if conversation exists for the given tenant and customer
	conv, err := s.convSvc.GetByTenantAndCustomer(tID, custID)
	if err != nil {
		return uuid.Nil, err
	}

	messageEntity := &domain.MessageEntity{
		ConversationID: conv.ID,
		SenderType:     senderType,
		Message:        message,
		CreatedAt:      time.Now(),
	}
	messageEntity.ConversationID = conv.ID
	return s.repo.Create(messageEntity)
}

func (s *MessageService) DeleteMessage(id uuid.UUID) error {
	return s.repo.Delete(id)
}
