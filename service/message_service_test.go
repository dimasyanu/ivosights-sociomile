package service

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/dimasyanu/ivosights-sociomile/config"
	"github.com/dimasyanu/ivosights-sociomile/constant"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra"
	"github.com/dimasyanu/ivosights-sociomile/internal/repository/mysqlrepo"
	"github.com/dimasyanu/ivosights-sociomile/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type MessageServiceTestSuite struct {
	svc         *MessageService
	tenantSvc   *TenantService
	rabbitMqCfg *infra.RabbitMQConfig
	mysqlCfg    *config.MysqlConfig
	db          *sql.DB
	mq          infra.QueueClient

	suite.Suite
}

func TestMessageServiceTestSuite(t *testing.T) {
	suite.Run(t, new(MessageServiceTestSuite))
}

func (s *MessageServiceTestSuite) SetupSuite() {
	const dbName = "test_messages"
	const envPath = "../.env"

	// Load configuration
	rabbitMqCfg := infra.NewRabbitMQConfig(envPath)
	s.mysqlCfg = config.NewMysqlConfig(envPath)
	s.rabbitMqCfg = infra.NewRabbitMQConfig(envPath)
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

	// Initialize repositories
	convRepo := mysqlrepo.NewConversationRepository(s.db)
	msgRepo := mysqlrepo.NewMessageRepository(s.db)
	tntRepo := mysqlrepo.NewTenantRepository(s.db)

	// Initialize RabbitMQ client
	s.mq, err = infra.NewRabbitMQClient(rabbitMqCfg)
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
	s.T().Logf("Dropping '%s' database ...", s.mysqlCfg.Database)
	if err := util.DropMysqlDatabase(s.mysqlCfg); err != nil {
		s.T().Fatalf("Failed to drop MySQL database: %v", err)
	}
}

// ===== Tests =====

func (s *MessageServiceTestSuite) TestGetMessages() {
}

func (s *MessageServiceTestSuite) TestCreateMessage() {
}

func (s *MessageServiceTestSuite) TestCreateMessageWithNewConversation_TriggersQueue() {
	// Create a new tenant
	tenant, err := s.tenantSvc.CreateTenant("Test Tenant")
	s.NoError(err)
	s.NotNil(tenant)

	// Random customer ID for testing
	customerID := uuid.New()

	// Attempt to create a message without an existing conversation
	msg := "Hello, I need help with my order."
	message, err := s.svc.CreateMessage(tenant.ID, customerID, constant.SenderTypeCustomer, msg)
	s.NoError(err)

	client, err := infra.NewRabbitMQClient(s.rabbitMqCfg)
	s.NoError(err)
	s.NotNil(client)

	published := client.GetPublishedMessages()
	s.Len(published, 1)

	expect := &infra.ConversationCreatedMessage{
		TenantID: tenant.ID,
		CustID:   customerID.String(),
		ConvID:   message.ConversationID.String(),
		Message:  msg,
	}
	expectBytes, err := json.Marshal(expect)
	s.NoError(err)
	s.Equal(string(expectBytes), string(published[0]))
}
