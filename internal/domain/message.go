package domain

import (
	"time"

	"github.com/google/uuid"
)

type MessageEntity struct {
	ID             uuid.UUID `sql:"column:id;primaryKey;autoIncrement"`
	ConversationID uuid.UUID `sql:"column:conversation_id;type:int;not null"`
	SenderType     string    `sql:"column:sender_type;type:varchar(20);not null"`
	Message        string    `sql:"column:message;type:text;not null"`
	CreatedAt      time.Time `sql:"column:created_at;autoCreateTime;not null"`
}

type Message struct {
	ID             uuid.UUID `json:"id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	SenderType     string    `json:"sender_type"`
	Message        string    `json:"message"`
	CreatedAt      time.Time `json:"created_at"`
}

func (m *MessageEntity) ToDto() *Message {
	return &Message{
		ID:             m.ID,
		ConversationID: m.ConversationID,
		SenderType:     m.SenderType,
		Message:        m.Message,
		CreatedAt:      m.CreatedAt,
	}
}
