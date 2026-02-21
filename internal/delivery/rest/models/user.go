package models

type UserCreateRequest struct {
	Name           string   `json:"name" validate:"required"`
	Email          string   `json:"email" validate:"required,email"`
	Roles          []string `json:"roles" validate:"required"`
	Password       string   `json:"password" validate:"required,min=6"`
	RepeatPassword string   `json:"repeat_password" validate:"required,eqfield=Password"`
}

type UserUpdateRequest struct {
	Name  *string  `json:"name,omitempty"`
	Email *string  `json:"email,omitempty" validate:"omitempty,email"`
	Roles []string `json:"roles,omitempty"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
	RepeatPassword  string `json:"repeat_password" validate:"required,eqfield=NewPassword"`
}
