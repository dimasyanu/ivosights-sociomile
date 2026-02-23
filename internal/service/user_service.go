package service

import (
	"strings"
	"time"

	"github.com/dimasyanu/ivosights-sociomile/internal/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/domain/repo"
	"github.com/dimasyanu/ivosights-sociomile/internal/utils"
	"github.com/google/uuid"
)

type UserService struct {
	repo repo.UserRepository
}

func NewUserService(repo repo.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUsers(filter *domain.UserFilter) (*domain.Paginated[domain.User], error) {
	entities, total, err := s.repo.GetList(filter)
	if err != nil {
		return nil, err
	}
	dtos := make([]domain.User, len(entities.Items))
	for i, entity := range entities.Items {
		dtos[i] = *entity.ToDto()
	}
	return &domain.Paginated[domain.User]{
		Items:    dtos,
		Page:     filter.Page,
		PageSize: filter.PageSize,
		Total:    total,
	}, nil
}

func (s *UserService) GetUserByID(id uuid.UUID) (*domain.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	dto := user.ToDto()
	return dto, nil
}

func (s *UserService) GetAvailableAgent() (*domain.User, error) {
	user, err := s.repo.GetAvailableAgent()
	if err != nil {
		return nil, err
	}
	return user.ToDto(), nil
}

func (s *UserService) GetUserByEmail(email string) (*domain.User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	dto := user.ToDto()
	return dto, nil
}

func (s *UserService) CreateUser(name, email, password string, roles []string, pic string) (uuid.UUID, error) {
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return uuid.Nil, err
	}
	user := &domain.UserEntity{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Roles:        strings.Join(roles, ","),
		CreatedAt:    time.Now(),
		CreatedBy:    pic,
		UpdatedAt:    time.Now(),
		UpdatedBy:    pic,
	}
	id, err := s.repo.Create(user)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (s *UserService) UpdateUser(id uuid.UUID, name *string, email *string, roles []string, pic string) (*domain.User, error) {
	user := &domain.UserEntity{
		ID:        id,
		UpdatedAt: time.Now(),
		UpdatedBy: pic,
	}
	data := map[string]any{}
	if name != nil {
		data["name"] = *name
	}
	if email != nil {
		data["email"] = *email
	}
	if len(roles) > 0 {
		data["roles"] = strings.Join(roles, ",")
	}

	user, err := s.repo.Update(user, data)
	if err != nil {
		return nil, err
	}
	return user.ToDto(), nil
}

func (s *UserService) UpdateUserPassword(id uuid.UUID, newPassword string, pic string) error {
	passwordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	user := &domain.UserEntity{
		ID:           id,
		PasswordHash: passwordHash,
		UpdatedAt:    time.Now(),
		UpdatedBy:    pic,
	}
	data := map[string]any{
		"password_hash": user.PasswordHash,
	}
	_, err = s.repo.Update(user, data)
	return err
}

func (s *UserService) DeleteUser(id uuid.UUID, pic string) error {
	now := time.Now()
	_, err := s.repo.Update(&domain.UserEntity{
		ID:        id,
		DeletedAt: &now,
	}, map[string]any{
		"deleted_at": &now,
	})
	return err
}
