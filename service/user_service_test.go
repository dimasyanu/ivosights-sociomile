package service

import (
	"database/sql"
	"log"
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
	log.SetFlags(log.Lshortfile)
	suite.Run(t, new(UserServiceTestSuite))
}

// Setup code before each test
func (s *UserServiceTestSuite) SetupSuite() {
	const dbName = "test_users"
	const envPath = "../.env"

	// Load configuration
	s.mysqlCfg = config.NewMysqlConfig(envPath)
	s.mysqlCfg.Database = dbName
	if err := util.CrateMysqlDatabase(s.mysqlCfg); err != nil {
		s.T().Fatalf("Failed to create MySQL database: %v", err)
	}

	// Initialize database and repositories
	var err error
	s.db, err = infra.NewMySQLDatabase(s.mysqlCfg)
	if err != nil {
		s.T().Fatalf("Failed to connect to database: %v", err)
	}
	s.repo = mysqlrepo.NewUserRepository(s.db)
	s.svc = NewUserService(s.repo)
}

// Cleanup code after each test
func (s *UserServiceTestSuite) TearDownSuite() {
	s.db.Close() // Close the database connection after all tests
	if err := util.DropMysqlDatabase(s.mysqlCfg); err != nil {
		s.T().Fatalf("Failed to drop MySQL database: %v", err)
	}
}

//========= Tests =========

func (s *UserServiceTestSuite) TestCreateUser() {
	const adminEmail = "admin@mail.com"
	r := &models.UserCreateRequest{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Roles:    []string{domain.RoleAgent},
		Password: "password123",
	}

	id, err := s.svc.CreateUser(r.Name, r.Email, r.Password, r.Roles, adminEmail)
	s.Require().NoError(err)
	s.Require().NotEqual(uuid.Nil, id)

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
}

func (s *UserServiceTestSuite) TestDeleteUser() {
}
