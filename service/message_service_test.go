package service

import (
	"database/sql"
	"testing"

	"github.com/dimasyanu/ivosights-sociomile/config"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra"
	"github.com/dimasyanu/ivosights-sociomile/internal/repository/mysqlrepo"
	"github.com/dimasyanu/ivosights-sociomile/util"
	"github.com/stretchr/testify/suite"
)

type MessageServiceTestSuite struct {
	svc      *MessageService
	mysqlCfg *config.MysqlConfig
	db       *sql.DB
	mq       infra.QueueClient

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

	// Initialize RabbitMQ client
	s.mq, err = infra.NewRabbitMQClient(rabbitMqCfg)
	if err != nil {
		s.T().Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	// Initialize services
	convSvc := NewConversationService(convRepo, s.mq)
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

func (s *MessageServiceTestSuite) TestGetMessages() {
}

func (s *MessageServiceTestSuite) TestCreateMessage() {
}

func (s *MessageServiceTestSuite) TestCreateMessageWithNewConversation_TriggersQueue() {

}
