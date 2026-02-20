package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/dimasyanu/ivosights-sociomile/common/httperr"
	"github.com/dimasyanu/ivosights-sociomile/config"
	"github.com/dimasyanu/ivosights-sociomile/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra"
	"github.com/dimasyanu/ivosights-sociomile/internal/repository/mysqlrepo"
	"github.com/dimasyanu/ivosights-sociomile/service"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	t *testing.T

	authSvc *service.AuthService
	db      infra.Database

	suite.Suite
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

// Setup code before each test
func (s *AuthTestSuite) SetupSuite() {
	// Load configuration
	mysqlCfg := config.NewMysqlConfig("../../.env")
	var err error
	s.db, err = infra.NewMySQLDatabase(mysqlCfg)
	if err != nil {
		s.T().Fatalf("Failed to connect to database: %v", err)
	}
	s.db.GetDb().Exec("DELETE FROM users;") // Clear users table before each test
	userRepo := mysqlrepo.NewUserRepository(s.db)
	s.authSvc = service.NewAuthService(userRepo)

	// Create a test user in the database
	password := "password!123"
	hashedPassword, err := infra.HashPassword(password)
	if err != nil {
		s.T().Fatalf("Failed to hash password: %v", err)
	}
	_, err = userRepo.CreateUser(&domain.UserEntity{
		Name:         "Test User",
		Email:        "test_login_success@mail.com",
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
		CreatedBy:    "system",
		UpdatedAt:    time.Now(),
		UpdatedBy:    "system",
	})
	if err != nil {
		s.T().Fatalf("Failed to create user: %v", err)
	}
}

// Cleanup code after each test
func (s *AuthTestSuite) TearDownSuite() {
	s.db.GetDb().Exec("DELETE FROM users;") // Clear users table after all tests
	s.db.Close()                            // Close the database connection after all tests
}

//========= Tests =========

// Test successful login
func (s *AuthTestSuite) TestLoginSuccess() {
	token, err := s.authSvc.Login("test_login_success@mail.com", "password!123")
	s.NoError(err, "Expected no error on successful login")
	s.NotEmpty(token, "Expected token to be generated on successful login")
}

// Test failed login
func (s *AuthTestSuite) TestLoginFailure() {
	token, err := s.authSvc.Login("test_login_failure@mail.com", "wrongpassword")
	s.Error(err)
	s.Equal(err.Error(), "Invalid email or password")
	s.True(errors.Is(err, httperr.Unauthorized), "Expected error to be Unauthorized")
	s.Empty(token, "Expected token to be empty on failed login")
}
