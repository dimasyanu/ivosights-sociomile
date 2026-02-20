package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserEntity struct {
	ID           uuid.UUID `gorm:"column:id;primaryKey;type:binary(16)"`
	Name         string    `gorm:"column:name;not null"`
	Email        string    `gorm:"column:email;unique;not null"`
	PasswordHash string    `gorm:"column:password_hash;not null"`

	CreatedAt time.Time       `gorm:"column:created_at;autoCreateTime;not null"`
	CreatedBy string          `gorm:"column:created_by;not null"`
	UpdatedAt time.Time       `gorm:"column:updated_at;autoUpdateTime;not null"`
	UpdatedBy string          `gorm:"column:updated_by;not null"`
	DeletedAt *gorm.DeletedAt `gorm:"column:deleted_at;nullable"`
}

type User struct {
	ID    uuid.UUID
	Name  string
	Email string
}
