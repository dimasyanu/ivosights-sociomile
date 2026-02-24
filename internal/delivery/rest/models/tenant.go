package models

type TenantCreateRequest struct {
	Name string `json:"name" validate:"required"`
}

type TenantUpdateRequest struct {
	Name string `json:"name" validate:"required"`
}
