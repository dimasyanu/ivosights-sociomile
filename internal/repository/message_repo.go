package repository

import (
	"github.com/dimasyanu/ivosights-sociomile/domain"
	"github.com/google/uuid"
)

type MessageRepository interface {
	GetByConversationID(id uuid.UUID, offset int64, limit int) ([]*domain.MessageEntity, int64, error)
	Create(message *domain.MessageEntity) (uuid.UUID, error)
	Delete(id uuid.UUID) error
}
