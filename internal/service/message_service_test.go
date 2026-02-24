package service

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/dimasyanu/ivosights-sociomile/config"
	"github.com/dimasyanu/ivosights-sociomile/internal/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/domain/constant"
	"github.com/dimasyanu/ivosights-sociomile/internal/domain/repo"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra/mysqlrepo"
	"github.com/dimasyanu/ivosights-sociomile/internal/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type MessageServiceTestSuite struct {
	svc       *MessageService
	tenantSvc *TenantService
	cfg       *config.Config
	convRepo  repo.ConversationRepository
	db        *sql.DB
	mq        infra.QueueClient

	suite.Suite
}

func TestMessageServiceTestSuite(t *testing.T) {
	suite.Run(t, new(MessageServiceTestSuite))
}

func (s *MessageServiceTestSuite) SetupSuite() {
	const dbName = "test_messages"
	const envPath = "../../.env"

	// Load configuration
	s.cfg = config.NewConfig(envPath)
	s.cfg.MySQL.Database = dbName

	// Create test database
	s.T().Logf("Creating database '%s'\n", s.cfg.MySQL.Database)
	if err := utils.CrateMysqlDatabase(envPath, s.cfg.MySQL); err != nil {
		s.T().Fatalf("Failed to create MySQL database: %v", err)
	}

	// Initialize database and repositories
	var err error
	s.db, err = infra.NewMySQLDatabase(s.cfg.MySQL)
	if err != nil {
		s.T().Fatalf("Failed to connect to database: %v", err)
	}
	s.T().Logf("Successfully connected to MySQL.")

	// Initialize repositories
	convRepo := mysqlrepo.NewConversationRepository(s.db)
	s.convRepo = convRepo
	msgRepo := mysqlrepo.NewMessageRepository(s.db)
	tntRepo := mysqlrepo.NewTenantRepository(s.db)

	// Initialize RabbitMQ client
	s.mq, err = infra.NewRabbitMQClient(s.cfg.RabbitMQ)
	if err != nil {
		s.T().Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	// Initialize services
	convSvc := NewConversationService(convRepo, s.mq)
	s.tenantSvc = NewTenantService(tntRepo)
	s.svc = NewMessageService(convSvc, msgRepo, s.mq)
}

func (s *MessageServiceTestSuite) TearDownSuite() {
	s.db.Close()
	s.mq.Close()
	s.T().Logf("Dropping '%s' database ...", s.cfg.MySQL.Database)
	if err := utils.DropMysqlDatabase(s.cfg.MySQL); err != nil {
		s.T().Fatalf("Failed to drop MySQL database: %v", err)
	}
}

func (s *MessageServiceTestSuite) TearDownTest() {
	s.mq.Clear()
}

// ===== Tests =====

func (s *MessageServiceTestSuite) TestGetMessages() {
}

func (s *MessageServiceTestSuite) TestCreateMessage() {
	// Create a new tenant
	tenant, err := s.tenantSvc.Create("Test Tenant")
	s.NoError(err)
	s.NotNil(tenant)

	// Random customer ID for testing
	customerID := uuid.New()

	// Set up a conversation
	conversation, err := s.convRepo.Create(&domain.ConversationEntity{
		ID:         uuid.New(),
		TenantID:   tenant.ID,
		CustomerID: customerID,
		Status:     constant.ConvStatusOpen,
		CreatedAt:  time.Now(),
	})
	s.NoError(err)
	s.NotNil(conversation)

	// Attempt to create a message without an existing conversation
	msg := "Hello, I need help with my order."
	_, err = s.svc.CreateMessage(tenant.ID, customerID, constant.SenderTypeCustomer, msg)
	s.NoError(err)

	client, err := infra.NewRabbitMQClient(s.cfg.RabbitMQ)
	s.NoError(err)
	s.NotNil(client)

	// Verify that no message was published to the queue
	published := client.GetPublishedMessages()
	s.Len(published, 0)
}

func (s *MessageServiceTestSuite) TestCreateMessageWithNewConversation_TriggersQueue() {
	// Create a new tenant
	tenant, err := s.tenantSvc.Create("Test Tenant")
	s.NoError(err)
	s.NotNil(tenant)

	// Random customer ID for testing
	customerID := uuid.New()

	// Attempt to create a message without an existing conversation
	msg := "Hello, I need help with my order."
	message, err := s.svc.CreateMessage(tenant.ID, customerID, constant.SenderTypeCustomer, msg)
	s.NoError(err)

	client, err := infra.NewRabbitMQClient(s.cfg.RabbitMQ)
	s.NoError(err)
	s.NotNil(client)

	// Verify that a message was published to the queue
	published := client.GetPublishedMessages()
	s.Len(published, 1)

	expect := &infra.ConversationCreatedMessage{
		TenantID: tenant.ID,
		CustID:   customerID,
		ConvID:   message.ConversationID,
		Message:  msg,
	}
	expectBytes, err := json.Marshal(expect)
	s.NoError(err)
	s.Equal(string(expectBytes), string(published[0]))
}

func (s *MessageServiceTestSuite) TestAssignAgentAfterConversationCreation() {
	// Create a new tenant
	tenant, err := s.tenantSvc.Create("Test Tenant")
	s.NoError(err)
	s.NotNil(tenant)

	// Random customer ID for testing
	customerID := uuid.New()

	// Attempt to create a message without an existing conversation
	msg := "Hello, I need help with my order."
	message, err := s.svc.CreateMessage(tenant.ID, customerID, constant.SenderTypeCustomer, msg)
	s.NoError(err)

	client, err := infra.NewRabbitMQClient(s.cfg.RabbitMQ)
	s.NoError(err)
	s.NotNil(client)

	conv, err := s.convRepo.GetByID(message.ConversationID)
	s.NoError(err)
	s.NotNil(conv)
	s.Nil(conv.AssignedAgentID)

}
