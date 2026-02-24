package models

type TenantCreateRequest struct {
	Name string `json:"name" validate:"required"`
}
