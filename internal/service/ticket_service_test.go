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

type TicketServiceTest struct {
	cfg *config.Config
	db  *sql.DB
	mq  infra.QueueClient

	repo     repo.TicketRepository
	convRepo repo.ConversationRepository
	msgRepo  repo.MessageRepository

	svc       *TicketService
	tenantSvc *TenantService

	suite.Suite
}

func TestTicketService(t *testing.T) {
	suite.Run(t, new(TicketServiceTest))
}

func (s *TicketServiceTest) SetupSuite() {
	const dbName = "test_tickets"
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

	// Initialize RabbitMQ client
	s.mq, err = infra.NewRabbitMQClient(s.cfg.RabbitMQ)
	if err != nil {
		s.T().Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	s.repo = mysqlrepo.NewTicketRepository(s.db)
	s.convRepo = mysqlrepo.NewConversationRepository(s.db)
	s.msgRepo = mysqlrepo.NewMessageRepository(s.db)

	convSvc := NewConversationService(s.convRepo, s.mq)
	s.svc = NewTicketService(s.repo, convSvc, s.mq)
	s.tenantSvc = NewTenantService(mysqlrepo.NewTenantRepository(s.db))
}

func (s *TicketServiceTest) TearDownSuite() {
	s.db.Close()
	s.T().Logf("Dropping '%s' database ...", s.cfg.MySQL.Database)
	if err := utils.DropMysqlDatabase(s.cfg.MySQL); err != nil {
		s.T().Fatalf("Failed to drop MySQL database: %v", err)
	}
	s.mq.Close()
}

// ===== Tests =====

func (s *TicketServiceTest) TestCreate() {
	ticket, err := s.svc.Create(uuid.New(), "New Ticket", "Ticket Description", 1, "admin@mail.com")
	s.Error(err) // Conversation should be exist

	tenant, err := s.tenantSvc.Create("Test Tenant")
	s.NoError(err)

	convId, err := s.convRepo.Create(&domain.ConversationEntity{
		ID:         uuid.New(),
		TenantID:   tenant.ID,
		CustomerID: uuid.New(),
		Status:     constant.ConvStatusAssigned,
		CreatedAt:  time.Now(),
	})
	s.NoError(err)
	s.NotNil(convId)

	ticket, err = s.svc.Create(convId, "New Ticket", "New ticket description", 1, "admin@mail.com")
	s.NoError(err)
	s.NotNil(ticket)
	s.Equal(constant.TicketStatusOpen, ticket.Status)

	// Another ticket for same conversation should return error
	ticket, err = s.svc.Create(convId, "Another Ticket", "Another ticket description", 1, "admin@mail.com")
	s.Error(err)
}

func (s *TicketServiceTest) TestGetByID() {
	tenant, err := s.tenantSvc.Create("Test Tenant")
	s.NoError(err)

	convId, err := s.convRepo.Create(&domain.ConversationEntity{
		ID:         uuid.New(),
		TenantID:   tenant.ID,
		CustomerID: uuid.New(),
		Status:     constant.ConvStatusAssigned,
		CreatedAt:  time.Now(),
	})
	s.NoError(err)
	s.NotNil(convId)

	ticket, err := s.svc.Create(convId, "New Ticket", "New ticket description", 1, "admin@mail.com")
	s.NoError(err)
	s.NotNil(ticket)
	s.Equal(constant.TicketStatusOpen, ticket.Status)

	fetched, err := s.svc.GetByID(ticket.ID)
	s.NoError(err)
	s.NotNil(fetched)
	s.Equal(ticket.ID, fetched.ID)
	s.Equal(ticket.ConversationID, fetched.ConversationID)
	s.Equal(ticket.Status, fetched.Status)
}

func (s *TicketServiceTest) TestGetByConversationID() {
	tenant, err := s.tenantSvc.Create("Test Tenant")
	s.NoError(err)

	convId, err := s.convRepo.Create(&domain.ConversationEntity{
		ID:         uuid.New(),
		TenantID:   tenant.ID,
		CustomerID: uuid.New(),
		Status:     constant.ConvStatusAssigned,
		CreatedAt:  time.Now(),
	})
	s.NoError(err)
	s.NotNil(convId)

	ticket, err := s.svc.Create(convId, "New Ticket", "New ticket description", 1, "admin@mail.com")
	s.NoError(err)
	s.NotNil(ticket)
	s.Equal(constant.TicketStatusOpen, ticket.Status)

	fetched, err := s.svc.GetByConversationID(convId)
	s.NoError(err)
	s.NotNil(fetched)
	s.Equal(ticket.ID, fetched.ID)
	s.Equal(ticket.ConversationID, fetched.ConversationID)
	s.Equal(ticket.Status, fetched.Status)
}

func (s *TicketServiceTest) TestUpdateStatus() {
	tenant, err := s.tenantSvc.Create("Test Tenant")
	s.NoError(err)

	convId, err := s.convRepo.Create(&domain.ConversationEntity{
		ID:         uuid.New(),
		TenantID:   tenant.ID,
		CustomerID: uuid.New(),
		Status:     constant.ConvStatusAssigned,
		CreatedAt:  time.Now(),
	})
	s.NoError(err)
	s.NotNil(convId)

	ticket, err := s.svc.Create(convId, "New Ticket", "New ticket description", 1, "admin@mail.com")
	s.NoError(err)
	s.NotNil(ticket)
	s.Equal(constant.TicketStatusOpen, ticket.Status)

	ticket, err = s.svc.UpdateStatus(ticket.ID, constant.TicketStatusInProgress)
	s.NoError(err)
	s.Equal(constant.TicketStatusInProgress, ticket.Status)

	fetched, err := s.svc.GetByID(ticket.ID)
	s.NoError(err)
	s.NotNil(fetched)
	s.Equal(constant.TicketStatusInProgress, fetched.Status)
}
