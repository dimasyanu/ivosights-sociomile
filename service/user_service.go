package service

import (
	"strings"
	"time"

	"github.com/dimasyanu/ivosights-sociomile/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/repository"
	"github.com/dimasyanu/ivosights-sociomile/util"
	"github.com/google/uuid"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUsers(filter *domain.UserFilter) (*domain.Paginated[domain.User], error) {
	entities, total, err := s.repo.GetUsers(filter)
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
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	dto := user.ToDto()
	return dto, nil
}

func (s *UserService) CreateUser(name, email, password string, roles []string, actor string) (uuid.UUID, error) {
	passwordHash, err := util.HashPassword(password)
	if err != nil {
		return uuid.Nil, err
	}
	user := &domain.UserEntity{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Roles:        strings.Join(roles, ","),
		CreatedAt:    time.Now(),
		CreatedBy:    actor,
		UpdatedAt:    time.Now(),
		UpdatedBy:    actor,
	}
	id, err := s.repo.CreateUser(user)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (s *UserService) UpdateUser(id uuid.UUID, name, email, password string) error {
	passwordHash, err := util.HashPassword(password)
	if err != nil {
		return err
	}
	user := &domain.UserEntity{
		ID:           id,
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
	}
	return s.repo.UpdateUser(user)
}

func (s *UserService) DeleteUser(id uuid.UUID) error {
	return s.repo.DeleteUser(id)
}
