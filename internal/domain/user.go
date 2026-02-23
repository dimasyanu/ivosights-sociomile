package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type UserEntity struct {
	ID           uuid.UUID `sql:"column:id;primaryKey;type:binary(16)"`
	Name         string    `sql:"column:name;not null"`
	Email        string    `sql:"column:email;unique;not null"`
	PasswordHash string    `sql:"column:password_hash;not null"`
	Roles        string    `sql:"column:roles;type:varchar(255);not null"` // Comma-separated role IDs

	CreatedAt time.Time  `sql:"column:created_at;autoCreateTime;not null"`
	CreatedBy string     `sql:"column:created_by;not null"`
	UpdatedAt time.Time  `sql:"column:updated_at;autoUpdateTime;not null"`
	UpdatedBy string     `sql:"column:updated_by;not null"`
	DeletedAt *time.Time `sql:"column:deleted_at;nullable"`
}

type User struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Roles []string  `json:"roles"`
}

type UserFilter struct {
	Name  *string `form:"name"`
	Email *string `form:"email"`

	Filter
}

func (u *UserEntity) ToDto() *User {
	return &User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Roles: strings.Split(u.Roles, ","),
	}
}
