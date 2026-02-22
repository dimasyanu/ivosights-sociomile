package mysqlrepo

import (
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
	panic("unimplemented")
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
