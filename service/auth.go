package service

import (
	"log"
	"strings"

	"github.com/dimasyanu/ivosights-sociomile/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra"
	"github.com/dimasyanu/ivosights-sociomile/internal/repository"
	"github.com/gofiber/fiber/v3"
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
	// Fetch user by email
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
		}
		log.Println(err.Error())
		return "", fiber.ErrInternalServerError
	}

	// Verify password
	if err = infra.CheckPasswordHash(password, user.PasswordHash); err != nil {
		if strings.Contains(err.Error(), "hashedPassword is not the hash") {
			return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
		}
		log.Println(err.Error())
		return "", fiber.ErrInternalServerError
	}

	// Generate JWT token
	token := "dummy_token" // Replace with actual token generation logic

	return token, nil
}

func (s *AuthService) Register(name, email, password string) error {
	_, err := s.userRepo.GetUserByEmail(email)
	if err == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Email already in use")
	}
	if !strings.Contains(err.Error(), "no rows in result set") {
		log.Fatal(err.Error())
		return fiber.ErrInternalServerError
	}

	passwordHash, err := infra.HashPassword(password)
	if err != nil {
		log.Fatal(err.Error())
		return fiber.ErrInternalServerError
	}

	user := &domain.UserEntity{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
	}

	if _, err = s.userRepo.CreateUser(user); err != nil {
		log.Fatal(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}
