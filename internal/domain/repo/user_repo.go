package repo

import (
	"github.com/dimasyanu/ivosights-sociomile/internal/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetList(filter *domain.UserFilter) (*domain.Paginated[domain.UserEntity], int64, error)
	GetAvailableAgent() (*domain.UserEntity, error)
	GetByID(id uuid.UUID) (*domain.UserEntity, error)
	GetByEmail(email string) (*domain.UserEntity, error)
	Create(user *domain.UserEntity) (uuid.UUID, error)
	Update(user *domain.UserEntity, data map[string]any) (*domain.UserEntity, error)
}
