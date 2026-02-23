package models

type ChannelPayload struct {
	TenantID   string `json:"tenant_id" validate:"required"`
	CustomerID string `json:"customer_id" validate:"required"`
	Message    string `json:"message" validate:"required"`
	Type       string `json:"type" validate:"required"`
}
