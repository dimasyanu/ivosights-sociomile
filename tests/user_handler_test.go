package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/dimasyanu/ivosights-sociomile/config"
	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest"
	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest/models"
	"github.com/dimasyanu/ivosights-sociomile/internal/domain"
	repository "github.com/dimasyanu/ivosights-sociomile/internal/domain/repo"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra/mysqlrepo"
	"github.com/dimasyanu/ivosights-sociomile/internal/service"
	"github.com/dimasyanu/ivosights-sociomile/internal/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

const envPath = "../.env"
const adminEmail = "admin@mail.com"

type UserHandlerTestSuite struct {
	mysqlCfg *config.MysqlConfig
	db       *sql.DB
	repo     repository.UserRepository
	svc      *service.UserService
	app      *fiber.App

	suite.Suite
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}

// Setup code before all tests
func (s *UserHandlerTestSuite) SetupSuite() {
	dbName := "test_user_handler"

	const envPath = "../.env"

	// Load configuration
	s.mysqlCfg = config.NewMysqlConfig(envPath)
	s.mysqlCfg.Database = dbName

	// Create test database
	s.T().Logf("Creating database '%s'\n", s.mysqlCfg.Database)
	if err := utils.CrateMysqlDatabase(envPath, s.mysqlCfg); err != nil {
		s.T().Fatalf("Failed to create MySQL database: %v", err)
	}

	// Initialize database and repositories
	var err error
	s.db, err = infra.NewMySQLDatabase(s.mysqlCfg)
	if err != nil {
		s.T().Fatalf("Failed to connect to database: %v", err)
	}
	s.T().Logf("Successfully connected to MySQL.")
	s.repo = mysqlrepo.NewUserRepository(s.db)
	s.svc = service.NewUserService(s.repo)

	s.app = s.MakeApp() // Initialize the Fiber app with routes
}

// Cleanup code after each test
func (s *UserHandlerTestSuite) TearDownTest() {
	s.db.Exec("DELETE FROM users WHERE email IS NOT ?", adminEmail) // Clean up users table after each test
}

// Cleanup code after all tests have run
func (s *UserHandlerTestSuite) TearDownSuite() {
	s.app.Shutdown() // Shutdown the Fiber app
	s.db.Close()     // Close the database connection after all tests

	s.T().Logf("Dropping '%s' database ...", s.mysqlCfg.Database)
	if err := utils.DropMysqlDatabase(s.mysqlCfg); err != nil {
		s.T().Fatalf("Failed to drop MySQL database: %v", err)
	}
}

// ====== Helper functions ======

func (s *UserHandlerTestSuite) MakeApp() *fiber.App {
	app := fiber.New()
	rest.RegisterRoutes(app, s.db, nil, envPath)

	return app
}

func login(s *suite.Suite, app *fiber.App) string {
	// Set up payload for login request
	payload := models.LoginRequest{
		Email:    adminEmail,
		Password: "my_secure_password",
	}
	payloadBytes, err := json.Marshal(payload)
	s.Require().NoError(err)
	jsonReader := bytes.NewReader(payloadBytes)

	// Create and send login request
	req := httptest.NewRequest("POST", "/api/v1/auth/login", jsonReader)
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	s.Require().NoError(err)
	s.Equal(200, res.StatusCode)

	// Read and parse response body
	resBytes, err := io.ReadAll(res.Body)
	s.Require().NoError(err)
	resBody := models.Res[models.LoginResponse]{}
	err = json.Unmarshal(resBytes, &resBody)
	s.Require().NoError(err)

	// Return the access token from the response
	return resBody.Data.AccessToken
}

// ====== Tests =====

func (s *UserHandlerTestSuite) TestAccessUser_Unauthorized() {
	req := httptest.NewRequest("GET", "/api/v1/backoffice/users/1", nil)
	res, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Equal(401, res.StatusCode)
}

func (s *UserHandlerTestSuite) TestCreateUser_Success() {
	// Login to get access token
	token := login(&s.Suite, s.app)

	// Set up payload for create user request
	payload := models.UserCreateRequest{
		Name:           "John Doe",
		Email:          "john.doe@mail.com",
		Roles:          []string{domain.RoleAgent},
		Password:       "password123",
		RepeatPassword: "password123",
	}
	payloadBytes, err := json.Marshal(payload)
	s.Require().NoError(err)
	jsonReader := bytes.NewReader(payloadBytes)

	// Create and send create user request with authorization header
	req := httptest.NewRequest("POST", "/api/v1/backoffice/users", jsonReader)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	res, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Equal(201, res.StatusCode)
}

func (s *UserHandlerTestSuite) TestCreateUser_DuplicateEmail() {
	const email = "existing.user@mail.com"
	s.svc.CreateUser("Existing User", email, "password123", []string{domain.RoleAgent}, adminEmail)

	// Login to get access token
	token := login(&s.Suite, s.app)

	// Set up payload for create user request
	payload := models.UserCreateRequest{
		Name:           "John Doe",
		Email:          email,
		Roles:          []string{domain.RoleAgent},
		Password:       "password123",
		RepeatPassword: "password123",
	}
	payloadBytes, err := json.Marshal(payload)
	s.Require().NoError(err)
	jsonReader := bytes.NewReader(payloadBytes)

	// Create and send create user request with authorization header
	req := httptest.NewRequest("POST", "/api/v1/backoffice/users", jsonReader)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	res, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Equal(400, res.StatusCode)
}

func (s *UserHandlerTestSuite) TestCreateUser_InvalidInput() {
	// Login to get access token
	token := login(&s.Suite, s.app)

	// Set up payload for create user request
	payload := models.UserCreateRequest{
		Name:           "John Doe",
		Email:          "john.doe@mail.com",
		Roles:          []string{domain.RoleAgent},
		Password:       "password123",
		RepeatPassword: "password1234", // Mismatched password
	}
	payloadBytes, err := json.Marshal(payload)
	s.Require().NoError(err)
	jsonReader := bytes.NewReader(payloadBytes)

	// Create and send create user request with authorization header
	req := httptest.NewRequest("POST", "/api/v1/backoffice/users", jsonReader)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	res, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Equal(400, res.StatusCode)
}

func (s *UserHandlerTestSuite) TestGetUserByID_Success() {
	id, err := s.svc.CreateUser("User", "user_by_id@mail.com", "password123", []string{domain.RoleAgent}, adminEmail)
	s.Require().NoError(err)

	token := login(&s.Suite, s.app)

	req := httptest.NewRequest("GET", "/api/v1/backoffice/users/"+id.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Equal(200, res.StatusCode)
}

func (s *UserHandlerTestSuite) TestGetUserByID_NotFound() {
	token := login(&s.Suite, s.app)

	req := httptest.NewRequest("GET", "/api/v1/backoffice/users/00000000-0000-0000-0000-000000000000", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Equal(404, res.StatusCode)
}

func (s *UserHandlerTestSuite) TestUpdateUser_Success() {
	// Create a user to update
	id, err := s.svc.CreateUser("User", "user_update@mail.com", "password123", []string{domain.RoleAgent}, adminEmail)
	s.Require().NoError(err)

	// Login to get access token
	token := login(&s.Suite, s.app)

	// Set up payload for update user request
	newName := "Updated User"
	newEmail := "updated_user@mail.com"
	payload := models.UserUpdateRequest{
		Name:  &newName,
		Email: &newEmail,
		Roles: []string{domain.RoleAgent},
	}
	payloadBytes, err := json.Marshal(payload)
	s.Require().NoError(err)
	jsonReader := bytes.NewReader(payloadBytes)

	// Create and send update user request with authorization header
	req := httptest.NewRequest("PUT", "/api/v1/backoffice/users/"+id.String(), jsonReader)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	res, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Equal(200, res.StatusCode)

	// Verify the user was updated in the database
	updatedUser, err := s.repo.GetByID(id)
	s.Require().NoError(err)
	s.Equal(newName, updatedUser.Name)
	s.Equal(newEmail, updatedUser.Email)
}

func (s *UserHandlerTestSuite) TestUpdateUser_NotFound() {
	s.svc.CreateUser("User", "update_unknown_user@mail.com", "password123", []string{domain.RoleAgent}, adminEmail)

	token := login(&s.Suite, s.app)

	// Set up payload for update user request
	newName := "Updated User"
	payload := models.UserUpdateRequest{
		Name:  &newName,
		Roles: []string{domain.RoleAgent},
	}
	payloadBytes, err := json.Marshal(payload)
	s.Require().NoError(err)
	jsonReader := bytes.NewReader(payloadBytes)

	// Create and send update user request with authorization header
	req := httptest.NewRequest("PUT", "/api/v1/backoffice/users/"+uuid.New().String(), jsonReader)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	res, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Equal(404, res.StatusCode)
}

func (s *UserHandlerTestSuite) TestDeleteUser_Success() {
	id, err := s.svc.CreateUser("User", "user_delete@mail.com", "password123", []string{domain.RoleAgent}, adminEmail)
	s.Require().NoError(err)

	token := login(&s.Suite, s.app)

	// Create and send delete user request with authorization header
	req := httptest.NewRequest("DELETE", "/api/v1/backoffice/users/"+id.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Equal(200, res.StatusCode)

	// Verify the user was deleted from the database
	user, err := s.repo.GetByID(id)
	s.Require().NoError(err)
	s.Require().NotNil(user.DeletedAt)
}

func (s *UserHandlerTestSuite) TestDeleteUser_NotFound() {
	token := login(&s.Suite, s.app)

	// Create and send delete user request with authorization header
	req := httptest.NewRequest("DELETE", "/api/v1/backoffice/users/"+uuid.New().String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Equal(404, res.StatusCode)
}
