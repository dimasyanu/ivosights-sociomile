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

func NewMessageRepository(db *sql.DB) repository.MessageRepository {
	return &MessageMySqlRepository{db: db}
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
	query := "DELETE FROM messages WHERE id = ?"
	_, err := m.db.ExecContext(context.Background(), query, id)
	return err
}

// GetByConversationID implements [repository.MessageRepository].
func (m *MessageMySqlRepository) GetByConversationID(id uuid.UUID, offset int64, limit int) ([]*domain.MessageEntity, int64, error) {
	query := "SELECT id, conversation_id, sender_type, message, created_at FROM messages WHERE conversation_id = ? ORDER BY created_at ASC LIMIT ? OFFSET ?"
	rows, err := m.db.QueryContext(context.Background(), query, id, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	messages := []*domain.MessageEntity{}
	for rows.Next() {
		var msg domain.MessageEntity
		if err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.SenderType, &msg.Message, &msg.CreatedAt); err != nil {
			return nil, 0, err
		}
		messages = append(messages, &msg)
	}

	var total int64
	countQuery := "SELECT COUNT(*) FROM messages WHERE conversation_id = ?"
	err = m.db.QueryRowContext(context.Background(), countQuery, id).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}
