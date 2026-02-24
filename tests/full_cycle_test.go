package tests

import (
	"database/sql"
	"testing"

	"github.com/dimasyanu/ivosights-sociomile/cmd/listener"
	"github.com/dimasyanu/ivosights-sociomile/config"
	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra"
	"github.com/dimasyanu/ivosights-sociomile/internal/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/suite"
)

const (
	agentEmail    = "agent@mail.com"
	agentPassword = "agent_password"
)

type FullCycleTestSuite struct {
	cfg *config.Config

	// Database
	db *sql.DB

	// Message Queue
	mq infra.QueueClient

	// Repositories

	// Services

	// App
	app *fiber.App

	// Worker
	worker *listener.QueueListener

	token string

	suite.Suite
}

func TestFullCycle(t *testing.T) {
	suite.Run(t, new(FullCycleTestSuite))
}

func (s *FullCycleTestSuite) makeApp() *fiber.App {
	app := fiber.New()
	rest.RegisterRoutes(app, s.db, s.mq, envPath)

	return app
}

func (s *FullCycleTestSuite) addWorker() {
	var err error
	s.worker, err = listener.NewQueueListener(s.cfg.RabbitMQ, s.db)
	if err != nil {
		s.T().Fatalf("Failed to initialize queue listener: %v", err)
	}

	// Start the worker
	go s.worker.Start(nil)
}

func (s *FullCycleTestSuite) SetupSuite() {
	dbName := "test_app_cycle"

	const envPath = "../.env"

	// Load configuration
	s.cfg = config.NewConfig(envPath)
	s.cfg.MySQL.Database = dbName

	// Create and initialize test database
	s.T().Logf("Creating database '%s'\n", s.cfg.MySQL.Database)
	if err := utils.CrateMysqlDatabase(envPath, s.cfg.MySQL); err != nil {
		s.T().Fatalf("Failed to create MySQL database: %v", err)
	}
	var err error
	s.db, err = infra.NewMySQLDatabase(s.cfg.MySQL)
	if err != nil {
		s.T().Fatalf("Failed to connect to database: %v", err)
	}
	s.T().Logf("Successfully connected to MySQL.")

	// Initialize the queue client
	s.mq, err = infra.NewRabbitMQClient(s.cfg.RabbitMQ)
	if err != nil {
		s.T().Fatalf("Failed to initialize RabbitMQ client: %v", err)
	}
	s.T().Logf("Successfully connected to RabbitMQ.")

	// Initialize the Fiber app with routes
	s.app = s.makeApp()

	// Start the worker
	s.addWorker()
}

func (s *FullCycleTestSuite) TearDownSuite() {
	// Drop test database
	s.T().Logf("Dropping database '%s'\n", s.cfg.MySQL.Database)
	if err := utils.DropMysqlDatabase(s.cfg.MySQL); err != nil {
		s.T().Fatalf("Failed to drop MySQL database: %v", err)
	}

	// Close database connection
	if err := s.db.Close(); err != nil {
		s.T().Logf("Failed to close database connection: %v", err)
	}
	s.T().Logf("Database connection closed successfully.")

	// Close RabbitMQ connection
	if err := s.mq.Close(); err != nil {
		s.T().Logf("Failed to close RabbitMQ connection: %v", err)
	}
	s.T().Logf("RabbitMQ connection closed successfully.")

	// Stop the worker
	s.worker.Close()
}

func (s *FullCycleTestSuite) TestAppCycle() {
	// Login as admin
	s.token = login(&s.Suite, s.app)

	// Create a user
	s.createAgentUser()

	// Create a tenant
	s.createTenant()

	// A message is received from the webhook
	// With an assumption that the tenantID is matched with the created tenant)
	s.webhookReceivedNewMessage()

	// Fetch the conversation

	// Login as agent

	// Fetch the conversation again and verify the message is there

	// Escalate the conversation to a ticket

	// Attempt to change the ticket status > SHOULD FAIL

	// Login as admin

	// Change the ticket status to in_progress > SHOULD SUCCEED
}
