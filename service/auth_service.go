package service

import (
	"database/sql"
	"errors"
	"log"

	"github.com/dimasyanu/ivosights-sociomile/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/repository"
	"github.com/dimasyanu/ivosights-sociomile/util"
	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo   repository.UserRepository
	JwtService *JwtService
}

func NewAuthService(userRepo repository.UserRepository, jwtService *JwtService) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		JwtService: jwtService,
	}
}

func (s *AuthService) Login(email, password string) (string, error) {
	// Fetch user by email
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
		}
		log.Println(err.Error())
		return "", fiber.ErrInternalServerError
	}

	// Verify password
	if err = util.CheckPasswordHash(password, user.PasswordHash); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
		}
		log.Println(err.Error())
		return "", fiber.ErrInternalServerError
	}

	// Generate JWT token
	token, err := s.JwtService.GenerateJWT(email)
	if err != nil {
		log.Println(err.Error())
		return "", fiber.ErrInternalServerError
	}

	return token, nil
}

func (s *AuthService) Register(name, email, password string) error {
	_, err := s.userRepo.GetUserByEmail(email)
	if err == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Email already in use")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	passwordHash, err := util.HashPassword(password)
	if err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	user := &domain.UserEntity{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
	}

	if _, err = s.userRepo.CreateUser(user); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *AuthService) ValidateToken(token string) (*domain.User, error) {
	email, err := s.JwtService.ValidateJWT(token)
	if err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrUnauthorized
	}

	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrUnauthorized
	}

	userDto := user.ToDto()
	return userDto, nil
}
