package models

import "github.com/google/uuid"

type ChannelPayload struct {
	TenantID   uint      `json:"tenant_id" validate:"required"`
	CustomerID uuid.UUID `json:"customer_id" validate:"required"`
	Message    string    `json:"message" validate:"required"`
	SenderType string    `json:"sender_type" validate:"required"`
}

type SendMessageRequest struct {
	Message string `json:"message" validate:"required"`
}
