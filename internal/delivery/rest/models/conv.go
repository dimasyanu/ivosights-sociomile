package models

type UpdateConversationStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=open closed assigned"`
}
