package models

type EscalateToTicketRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Priority    int8   `json:"priority" validate:"required,min=1,max=5"`
}

type UpdateTicketStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=open in_progress closed"`
}
