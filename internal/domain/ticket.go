package domain

import (
	"time"

	"github.com/google/uuid"
)

type TicketEntity struct {
	ID             uuid.UUID  `json:"id"`
	ConversationID uuid.UUID  `json:"conversation_id"`
	Title          string     `json:"subject"`
	Description    string     `json:"description"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	CreatedBy      string     `json:"created_by"`
	UpdatedAt      time.Time  `json:"updated_at"`
	UpdatedBy      string     `json:"updated_by"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	DeletedBy      string     `json:"deleted_by,omitempty"`
}

type Ticket struct {
	ID             uuid.UUID  `json:"id"`
	ConversationID uuid.UUID  `json:"conversation_id"`
	Subject        string     `json:"subject"`
	Description    string     `json:"description"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	CreatedBy      string     `json:"created_by"`
	UpdatedAt      time.Time  `json:"updated_at"`
	UpdatedBy      string     `json:"updated_by"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	DeletedBy      string     `json:"deleted_by,omitempty"`
}

type TicketFilter struct {
	ConversationID *uuid.UUID `json:"conversation_id,omitempty"`
	Status         *string    `json:"status,omitempty"`
}

func (e *TicketEntity) ToDto() *Ticket {
	return &Ticket{
		ID:             e.ID,
		ConversationID: e.ConversationID,
		Subject:        e.Title,
		Description:    e.Description,
		Status:         e.Status,
		CreatedAt:      e.CreatedAt,
		CreatedBy:      e.CreatedBy,
		UpdatedAt:      e.UpdatedAt,
		UpdatedBy:      e.UpdatedBy,
		DeletedAt:      e.DeletedAt,
		DeletedBy:      e.DeletedBy,
	}
}
