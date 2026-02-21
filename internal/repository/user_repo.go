package repository

import (
	"github.com/dimasyanu/ivosights-sociomile/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetUsers(filter *domain.UserFilter) (*domain.Paginated[domain.UserEntity], int64, error)
	GetUserByID(id uuid.UUID) (*domain.UserEntity, error)
	GetUserByEmail(email string) (*domain.UserEntity, error)
	CreateUser(user *domain.UserEntity) (uuid.UUID, error)
	UpdateUser(user *domain.UserEntity, data map[string]any) (*domain.UserEntity, error)
}
