package service

import (
	"database/sql"
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

type ConversationServiceTestSuite struct {
	mysqlCfg    *config.MysqlConfig
	rabbitMqCfg *config.RabbitMQConfig
	db          *sql.DB

	mq        infra.QueueClient
	repo      repo.ConversationRepository
	tenantSvc *TenantService
	svc       *ConversationService

	suite.Suite
}

func TestConversationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ConversationServiceTestSuite))
}

func (s *ConversationServiceTestSuite) SetupSuite() {
	const dbName = "test_conversations"
	const envPath = "../../.env"

	// Load configuration
	s.mysqlCfg = config.NewMysqlConfig(envPath)
	s.rabbitMqCfg = config.NewRabbitMQConfig(envPath)
	s.mysqlCfg.Database = dbName

	// Initialize RabbitMQ
	var err error
	s.mq, err = infra.NewRabbitMQClient(s.rabbitMqCfg)

	// Create test database
	s.T().Logf("Creating database '%s'\n", s.mysqlCfg.Database)
	if err := utils.CrateMysqlDatabase(envPath, s.mysqlCfg); err != nil {
		s.T().Fatalf("Failed to create MySQL database: %v", err)
	}

	// Initialize database and repositories
	s.db, err = infra.NewMySQLDatabase(s.mysqlCfg)
	if err != nil {
		s.T().Fatalf("Failed to connect to database: %v", err)
	}
	s.T().Logf("Successfully connected to MySQL.")

	tenantRepo := mysqlrepo.NewTenantRepository(s.db)
	s.repo = mysqlrepo.NewConversationRepository(s.db)

	s.tenantSvc = NewTenantService(tenantRepo)
	s.svc = NewConversationService(s.repo, s.mq)
}

func (s *ConversationServiceTestSuite) TearDownSuite() {
	s.mq.Close()
	s.db.Close()

	s.T().Logf("Dropping '%s' database ...", s.mysqlCfg.Database)
	if err := utils.DropMysqlDatabase(s.mysqlCfg); err != nil {
		s.T().Fatalf("Failed to drop MySQL database: %v", err)
	}
}

func (s *ConversationServiceTestSuite) TestGetByID() {
	tenant, err := s.tenantSvc.Create("Test Tenant")
	s.Require().NoError(err)

	custID := uuid.New()

	convID, err := s.repo.Create(&domain.ConversationEntity{
		TenantID:   tenant.ID,
		CustomerID: custID,
		Status:     constant.ConvStatusOpen,
		CreatedAt:  time.Now(),
	})
	if err != nil {
		s.T().Logf("An error occured: %v\n", err)
	}

	conv, err := s.svc.GetByID(convID)
	s.Require().NoError(err)
	s.Require().Equal(convID, conv.ID)
	s.NotNil(conv.ID)
	s.NotEqual(uuid.Nil, conv.ID)
}
