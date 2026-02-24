package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dimasyanu/ivosights-sociomile/cmd/listener"
	"github.com/dimasyanu/ivosights-sociomile/config"
	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest"
	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest/models"
	"github.com/dimasyanu/ivosights-sociomile/internal/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/domain/constant"
	repository "github.com/dimasyanu/ivosights-sociomile/internal/domain/repo"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra/mysqlrepo"
	"github.com/dimasyanu/ivosights-sociomile/internal/service"
	"github.com/dimasyanu/ivosights-sociomile/internal/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type MessageHandlerTestSuite struct {
	mysqlCfg    *config.MysqlConfig
	rabbitMqCfg *config.RabbitMQConfig

	db       *sql.DB
	userRepo repository.UserRepository
	convRepo repository.ConversationRepository

	svc       *service.UserService
	tenantSvc *service.TenantService
	mq        infra.QueueClient

	app *fiber.App

	suite.Suite
}

func TestMessageHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(MessageHandlerTestSuite))
}

func (s *MessageHandlerTestSuite) SetupSuite() {
	dbName := "test_message_handler"

	const envPath = "../.env"

	// Load configuration
	s.mysqlCfg = config.NewMysqlConfig(envPath)
	s.mysqlCfg.Database = dbName
	s.rabbitMqCfg = config.NewRabbitMQConfig(envPath)

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

	s.userRepo = mysqlrepo.NewUserRepository(s.db)
	s.convRepo = mysqlrepo.NewConversationRepository(s.db)
	tenantRepo := mysqlrepo.NewTenantRepository(s.db)

	s.svc = service.NewUserService(s.userRepo)
	s.tenantSvc = service.NewTenantService(tenantRepo)

	s.mq, err = infra.NewRabbitMQClient(s.rabbitMqCfg) // Initialize the queue client
	if err != nil {
		s.T().Fatalf("Failed to initialize RabbitMQ client: %v", err)
	}

	s.app = s.MakeApp() // Initialize the Fiber app with routes
}

func (s *MessageHandlerTestSuite) TearDownSuite() {
	// Drop test database
	s.T().Logf("Dropping database '%s'\n", s.mysqlCfg.Database)
	if err := utils.DropMysqlDatabase(s.mysqlCfg); err != nil {
		s.T().Fatalf("Failed to drop MySQL database: %v", err)
	}

	// Close database connection
	if err := s.db.Close(); err != nil {
		s.T().Logf("Failed to close database connection: %v", err)
	} else {
		s.T().Logf("Database connection closed successfully.")
	}

	// Close RabbitMQ connection
	if err := s.mq.Close(); err != nil {
		s.T().Logf("Failed to close RabbitMQ connection: %v", err)
	} else {
		s.T().Logf("RabbitMQ connection closed successfully.")
	}
}

func (s *MessageHandlerTestSuite) TearDownTest() {
	// Clear database tables after each test to ensure isolation
	s.db.Exec("DELETE FROM users")

	// Clear the queue after each test to ensure isolation
	if err := s.mq.Clear(); err != nil {
		s.T().Logf("Failed to clear RabbitMQ queue: %v", err)
	} else {
		s.T().Logf("RabbitMQ queue cleared successfully.")
	}
}

// ====== Helper functions ======

func (s *MessageHandlerTestSuite) MakeApp() *fiber.App {
	app := fiber.New()
	rest.RegisterRoutes(app, s.db, s.mq, envPath)

	return app
}

// ====== Tests ======

func (s *MessageHandlerTestSuite) TestHandleMessageCreatedWithNewConversation() {
	// Set up tenant
	tenant, err := s.tenantSvc.Create("Test Tenant")
	s.Require().NoError(err)

	// Set up an available agent
	pass, _ := utils.HashPassword("password123")
	agent := &domain.UserEntity{
		Name:         "Smith",
		Email:        "smith.agent@example.com",
		Roles:        "agent",
		PasswordHash: pass,
		CreatedAt:    time.Now(),
		CreatedBy:    "system",
		UpdatedAt:    time.Now(),
		UpdatedBy:    "system",
	}
	agentID, err := s.userRepo.Create(agent)
	s.Require().NoError(err)
	agent.ID = agentID

	// Prepare request payload
	payload := &models.ChannelPayload{
		TenantID:   tenant.ID,
		CustomerID: uuid.New(),
		Message:    "Hello, I need help with my order.",
		SenderType: constant.SenderTypeCustomer,
	}
	payloadBytes, err := json.Marshal(payload)
	s.Require().NoError(err)
	jsonReader := bytes.NewReader(payloadBytes)

	// Create and send create user request with authorization header
	req := httptest.NewRequest("POST", "/api/v1/channel/webhook", jsonReader)
	req.Header.Set("Content-Type", "application/json")
	res, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Equal(201, res.StatusCode)

	// Verify that a new conversation is not assigned to any agent
	conv, err := s.convRepo.GetByTenantAndCustomer(tenant.ID, payload.CustomerID)
	s.Require().NoError(err)
	s.Equal(tenant.ID, conv.TenantID)
	s.Equal(payload.CustomerID, conv.CustomerID)
	s.Nil(conv.AssignedAgentID)

	// Start the worker
	listener, err := listener.NewQueueListener(s.rabbitMqCfg, s.db)
	if err != nil {
		s.T().Fatalf("Failed to initialize queue listener: %v", err)
	}
	defer listener.Close()

	go listener.Start(nil)

	// // Wait for the worker to process the message
	time.Sleep(time.Millisecond * 200)

	// Verify that the the conversation is now assigned to an agent
	conv, err = s.convRepo.GetByTenantAndCustomer(tenant.ID, payload.CustomerID)
	s.NoError(err)
	s.NotNil(conv.AssignedAgentID)
	s.Equal(agent.ID, *conv.AssignedAgentID)
}
