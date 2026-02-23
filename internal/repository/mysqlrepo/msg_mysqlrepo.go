package mysqlrepo

import (
	"context"
	"database/sql"

	"github.com/dimasyanu/ivosights-sociomile/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/repository"
	"github.com/google/uuid"
)

type MessageMySqlRepository struct {
	db *sql.DB
}

// Create implements [repository.MessageRepository].
func (m *MessageMySqlRepository) Create(message *domain.MessageEntity) (uuid.UUID, error) {
	id := uuid.New()
	pairs := map[string]any{
		"id":              id,
		"conversation_id": message.ConversationID,
		"sender_type":     message.SenderType,
		"message":         message.Message,
		"created_at":      message.CreatedAt,
	}
	cols, slots, vals := MapForCreate(pairs)

	query := "INSERT INTO messages (" + cols + ") VALUES (" + slots + ")"
	_, err := m.db.ExecContext(context.Background(), query, vals...)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

// Delete implements [repository.MessageRepository].
func (m *MessageMySqlRepository) Delete(id uuid.UUID) error {
	panic("unimplemented")
}

// GetByConversationID implements [repository.MessageRepository].
func (m *MessageMySqlRepository) GetByConversationID(id uuid.UUID, offset int64, limit int) ([]*domain.MessageEntity, int64, error) {
	panic("unimplemented")
}

func NewMessageRepository(db *sql.DB) repository.MessageRepository {
	return &MessageMySqlRepository{db: db}
}
