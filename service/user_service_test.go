package service

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/dimasyanu/ivosights-sociomile/config"
	"github.com/dimasyanu/ivosights-sociomile/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest/models"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra"
	"github.com/dimasyanu/ivosights-sociomile/internal/repository"
	"github.com/dimasyanu/ivosights-sociomile/internal/repository/mysqlrepo"
	"github.com/dimasyanu/ivosights-sociomile/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type UserServiceTestSuite struct {
	t *testing.T

	mysqlCfg *config.MysqlConfig
	db       *sql.DB
	repo     repository.UserRepository
	svc      *UserService

	suite.Suite
}

func TestUserServiceTestSuite(t *testing.T) {
	// log.SetFlags(log.Lshortfile)
	suite.Run(t, new(UserServiceTestSuite))
}

// Setup code before each test
func (s *UserServiceTestSuite) SetupSuite() {
	const dbName = "test_users"
	const envPath = "../.env"

	// Load configuration
	s.mysqlCfg = config.NewMysqlConfig(envPath)
	s.mysqlCfg.Database = dbName

	// Create test database
	s.T().Logf("Creating database '%s'\n", s.mysqlCfg.Database)
	if err := util.CrateMysqlDatabase(s.mysqlCfg); err != nil {
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
	s.svc = NewUserService(s.repo)
}

// Cleanup code after each test
func (s *UserServiceTestSuite) TearDownSuite() {
	s.db.Close() // Close the database connection after all tests

	s.T().Logf("Dropping '%s' database ...", s.mysqlCfg.Database)
	if err := util.DropMysqlDatabase(s.mysqlCfg); err != nil {
		s.T().Fatalf("Failed to drop MySQL database: %v", err)
	}
}

func RemoveUserById(db *sql.DB, id uuid.UUID) error {
	query := "DELETE FROM users WHERE id = UUID_TO_BIN(?)"
	_, err := db.Exec(query, id.String())
	return err
}

//========= Tests =========

func (s *UserServiceTestSuite) TestCreateUser() {
	const adminEmail = "admin@mail.com"
	r := &models.UserCreateRequest{
		Name:     "John Doe",
		Email:    "john.doe1@example.com",
		Roles:    []string{domain.RoleAgent},
		Password: "password123",
	}

	id, err := s.svc.CreateUser(r.Name, r.Email, r.Password, r.Roles, adminEmail)
	s.Require().NoError(err)
	s.Require().NotEqual(uuid.Nil, id)
	defer RemoveUserById(s.db, id)

	// Verify user was created in the database
	user, err := s.repo.GetUserByID(id)
	s.Require().NoError(err)
	s.Require().NotNil(user)
	s.Equal(r.Name, user.Name)
	s.Equal(r.Email, user.Email)
	s.Require().NoError(util.CheckPasswordHash(r.Password, user.PasswordHash))
	s.Equal(adminEmail, user.CreatedBy)
	s.Equal(adminEmail, user.UpdatedBy)
	s.Equal(strings.Join(r.Roles, ","), user.Roles)
}

func (s *UserServiceTestSuite) TestGetUsers() {
	filter := &domain.UserFilter{
		Filter: domain.Filter{
			Page:     1,
			PageSize: 25,
		},
	}

	paginated, err := s.svc.GetUsers(filter)
	s.Require().NoError(err)
	s.Require().NotNil(paginated)
	s.Equal(1, paginated.Page)
	s.Equal(25, paginated.PageSize)
}

func (s *UserServiceTestSuite) TestUpdateUser() {
	const adminEmail = "admin@mail.com"
	const name = "John Doe"
	const email = "john.doe2@example.com"
	const password = "password123"
	var roles = []string{domain.RoleAgent}

	id, err := s.svc.CreateUser(name, email, password, roles, adminEmail)
	s.Require().NoError(err)
	defer RemoveUserById(s.db, id)

	updatedName := "John Doe Updated"
	updateReq := &models.UserUpdateRequest{
		Name: &updatedName,
	}
	_, err = s.svc.UpdateUser(id, updateReq.Name, nil, nil, adminEmail)
	s.Require().NoError(err)

	// Verify user was updated in the database
	user, err := s.repo.GetUserByID(id)
	s.Require().NoError(err)
	s.Require().NotNil(user)
	s.Equal(updatedName, user.Name)
	s.Equal(adminEmail, user.UpdatedBy)

	// Update user roles
	updatedRoles := []string{domain.RoleAdmin}
	updateReq = &models.UserUpdateRequest{
		Roles: updatedRoles,
	}
	_, err = s.svc.UpdateUser(id, nil, nil, updateReq.Roles, adminEmail)
	s.Require().NoError(err)

	// Verify user roles were updated in the database
	user, err = s.repo.GetUserByID(id)
	s.Require().NoError(err)
	s.Require().NotNil(user)
	s.Equal(strings.Join(updatedRoles, ","), user.Roles)
	s.Equal(adminEmail, user.UpdatedBy)
}

func (s *UserServiceTestSuite) TestDeleteUser() {
	const adminEmail = "admin@mail.com"
	const name = "John Doe"
	const email = "john.doe3@example.com"
	const password = "password123"
	var roles = []string{domain.RoleAgent}

	id, err := s.svc.CreateUser(name, email, password, roles, adminEmail)
	s.Require().NoError(err)
	defer RemoveUserById(s.db, id)

	err = s.svc.DeleteUser(id, adminEmail)
	s.Require().NoError(err)

	// Verify user was deleted from the database
	user, err := s.repo.GetUserByID(id)
	s.Require().NoError(err)
	s.Require().NotNil(user)
	s.Require().NotNil(user.DeletedAt)
}
