package service

import (
	"github.com/dimasyanu/ivosights-sociomile/common/httperr"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra"
	"github.com/dimasyanu/ivosights-sociomile/internal/repository"
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", httperr.Unauthorized
	}

	if err = infra.CheckPasswordHash(password, user.PasswordHash); err != nil {
		return "", httperr.Unauthorized
	}

	// Generate JWT token or any other authentication token here
	token := "dummy_token" // Replace with actual token generation logic

	return token, nil
}
