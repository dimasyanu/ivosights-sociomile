package service

import (
	"database/sql"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/dimasyanu/ivosights-sociomile/config"
	"github.com/dimasyanu/ivosights-sociomile/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra"
	"github.com/dimasyanu/ivosights-sociomile/internal/repository/mysqlrepo"
	"github.com/dimasyanu/ivosights-sociomile/util"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/suite"
)

type AuthServiceTestSuite struct {
	t *testing.T

	mysqlCfg *config.MysqlConfig
	authSvc  *AuthService
	db       *sql.DB

	suite.Suite
}

func TestAuthTestSuite(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	suite.Run(t, new(AuthServiceTestSuite))
}

// Setup code before each test
func (s *AuthServiceTestSuite) SetupSuite() {
	const dbName = "test_login"
	const envPath = "../.env"

	var err error

	// Load configuration
	s.mysqlCfg = config.NewMysqlConfig(envPath)
	s.mysqlCfg.Database = dbName
	if err = util.CrateMysqlDatabase(s.mysqlCfg); err != nil {
		s.T().Fatalf("Failed to create MySQL database: %v", err)
	}

	// Initialize database and repositories
	s.db, err = infra.NewMySQLDatabase(s.mysqlCfg)
	if err != nil {
		s.T().Fatalf("Failed to connect to database: %v", err)
	}
	userRepo := mysqlrepo.NewUserRepository(s.db)
	jwtService := NewJwtService(config.NewJwtConfig(envPath))
	s.authSvc = NewAuthService(userRepo, jwtService)

	// Create a test user in the database
	password := "password!123"
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		s.T().Fatalf("Failed to hash password: %v", err)
	}
	_, err = userRepo.CreateUser(&domain.UserEntity{
		Name:         "Test User",
		Email:        "test_login@mail.com",
		PasswordHash: hashedPassword,
		Roles:        domain.RoleAgent,
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
func (s *AuthServiceTestSuite) TearDownSuite() {
	s.db.Exec("DELETE FROM users;") // Clear users table after all tests
	s.db.Close()                    // Close the database connection after all tests
	if err := util.DropMysqlDatabase(s.mysqlCfg); err != nil {
		s.T().Fatalf("Failed to drop MySQL database: %v", err)
	}
}

//========= Tests =========

// Test successful login
func (s *AuthServiceTestSuite) TestLoginSuccess() {
	token, err := s.authSvc.Login("test_login@mail.com", "password!123")
	s.NoError(err, "Expected no error on successful login")
	s.NotEmpty(token, "Expected token to be generated on successful login")
}

// Test failed login
func (s *AuthServiceTestSuite) TestLoginFailure() {
	token, err := s.authSvc.Login("test_wrong_email@mail.com", "password!123")
	s.Error(err)
	s.Equal("Invalid email or password", err.Error())
	s.Empty(token, "Expected token to be empty on failed login")

	ferr := &fiber.Error{}
	if errors.As(err, &ferr) {
		s.Equal(fiber.StatusUnauthorized, ferr.Code, "Expected status code to be 401 Unauthorized")
		s.Equal("Invalid email or password", ferr.Message, "Expected error message to be 'Invalid email or password'")
	}

	token, err = s.authSvc.Login("test_login@mail.com", "wrongpassword")
	s.Error(err)
	s.Equal("Invalid email or password", err.Error())
	s.Empty(token, "Expected token to be empty on failed login")
	if errors.As(err, &ferr) {
		s.Equal(fiber.StatusUnauthorized, ferr.Code, "Expected status code to be 401 Unauthorized")
		s.Equal("Invalid email or password", ferr.Message, "Expected error message to be 'Invalid email or password'")
	}
}
