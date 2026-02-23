package trinorepo

import (
	"database/sql"

	"github.com/dimasyanu/ivosights-sociomile/internal/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/domain/repo"
	"github.com/google/uuid"
)

// TODO: Implement repo for user domain using Trino, MinIO, and Iceberg.

type UserTrinoRepository struct {
	db *sql.DB
}

func NewUserTrinoRepository(db *sql.DB) repo.UserRepository {
	return &UserTrinoRepository{db: db}
}

// Create implements [repo.UserRepository].
func (u *UserTrinoRepository) Create(user *domain.UserEntity) (uuid.UUID, error) {
	panic("unimplemented")
}

// GetAvailableAgent implements [repo.UserRepository].
func (u *UserTrinoRepository) GetAvailableAgent() (*domain.UserEntity, error) {
	panic("unimplemented")
}

// GetByEmail implements [repo.UserRepository].
func (u *UserTrinoRepository) GetByEmail(email string) (*domain.UserEntity, error) {
	panic("unimplemented")
}

// GetByID implements [repo.UserRepository].
func (u *UserTrinoRepository) GetByID(id uuid.UUID) (*domain.UserEntity, error) {
	panic("unimplemented")
}

// GetList implements [repo.UserRepository].
func (u *UserTrinoRepository) GetList(filter *domain.UserFilter) (*domain.Paginated[domain.UserEntity], int64, error) {
	panic("unimplemented")
}

// Update implements [repo.UserRepository].
func (u *UserTrinoRepository) Update(user *domain.UserEntity, data map[string]any) (*domain.UserEntity, error) {
	panic("unimplemented")
}
