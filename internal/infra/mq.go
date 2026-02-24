package infra

import "github.com/google/uuid"

type QueueClient interface {
	Publish(queue string, message []byte) error
	Consume(queue string) (<-chan []byte, error)
	GetPublishedMessages() [][]byte
	Clear() error
	Close() error
}

type ConversationCreatedMessage struct {
	TenantID uint      `json:"tenant_id"`
	CustID   uuid.UUID `json:"cust_id"`
	ConvID   uuid.UUID `json:"conv_id"`
	Message  string    `json:"message"`
}

type ConversationStatusUpdatedMessage struct {
	ConvID uuid.UUID `json:"conv_id"`
	Status string    `json:"status"`
}

type TicketCreatedMessage struct {
	ConvID   uuid.UUID `json:"conv_id"`
	TicketID uuid.UUID `json:"ticket_id"`
	Status   string    `json:"status"`
}
