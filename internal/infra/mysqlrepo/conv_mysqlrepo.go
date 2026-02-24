package mysqlrepo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/dimasyanu/ivosights-sociomile/internal/domain"
	repository "github.com/dimasyanu/ivosights-sociomile/internal/domain/repo"
	"github.com/google/uuid"
)

type ConversationMySqlRepository struct {
	db *sql.DB
}

func NewConversationRepository(db *sql.DB) repository.ConversationRepository {
	// Initialize your database connection here and assign it to the db field
	return &ConversationMySqlRepository{
		db: db, // Replace with actual DB connection
	}
}

func (r *ConversationMySqlRepository) GetList(f *domain.ConversationFilter) ([]domain.ConversationEntity, uint64, error) {
	// Build the SQL baseQuery based on the filter parameters
	var sb strings.Builder
	sb.WriteString("SELECT %s FROM conversations WHERE 1=1")

	args := []any{}

	if f.TenantID != nil {
		sb.WriteString(" AND tenant_id = ?")
		args = append(args, *f.TenantID)
	}
	if f.CustomerID != nil {
		sb.WriteString(" AND customer_id = UUID_TO_BIN(?)")
		args = append(args, f.CustomerID.String())
	}
	if f.AssignedAgentID != nil {
		sb.WriteString(" AND assigned_agent_id = UUID_TO_BIN(?)")
		args = append(args, f.AssignedAgentID.String())
	}
	if f.Status != nil {
		sb.WriteString(" AND status = ?")
		args = append(args, *f.Status)
	}

	limit := f.PageSize
	offset := (f.Page - 1) * f.PageSize
	query := fmt.Sprintf(sb.String(), "id, tenant_id, customer_id, assigned_agent_id, status") + " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var conversations []domain.ConversationEntity
	for rows.Next() {
		var conv domain.ConversationEntity
		err := rows.Scan(&conv.ID, &conv.TenantID, &conv.CustomerID, &conv.AssignedAgentID, &conv.Status)
		if err != nil {
			return nil, 0, err
		}
		conversations = append(conversations, conv)
	}

	// Get total count for pagination
	var total uint64
	countQuery := fmt.Sprintf(sb.String(), "COUNT(*) AS total")

	err = r.db.QueryRowContext(context.Background(), countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return conversations, total, nil
}

// Create implements [repository.ConversationRepository].
func (r *ConversationMySqlRepository) Create(c *domain.ConversationEntity) (uuid.UUID, error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	pairs := map[string]any{
		"id":          c.ID,
		"tenant_id":   c.TenantID,
		"customer_id": c.CustomerID,
		"status":      c.Status,
		"created_at":  c.CreatedAt,
	}
	cols, slots, vals := MapForCreate(pairs)

	query := "INSERT INTO conversations (" + cols + ") VALUES (" + slots + ")"
	_, err := r.db.ExecContext(context.Background(), query, vals...)
	if err != nil {
		return uuid.Nil, err
	}

	return c.ID, nil
}

// GetByID implements [repository.ConversationRepository].
func (r *ConversationMySqlRepository) GetByID(id uuid.UUID) (*domain.ConversationEntity, error) {
	conversation := &domain.ConversationEntity{}
	query := "SELECT %s FROM conversations c %s WHERE c.id = UUID_TO_BIN(?)"
	selects := []string{
		"c.id",
		"c.tenant_id",
		"c.customer_id",
		"c.assigned_agent_id",
		"c.status",
		"t.name AS tenant_name",
		"a.name AS agent_name",
		"a.email AS agent_email",
	}
	join := []string{
		"LEFT JOIN tenants t ON c.tenant_id = t.id",
		"LEFT JOIN users a ON c.assigned_agent_id = a.id",
	}
	query = fmt.Sprintf(query, strings.Join(selects, ", "), strings.Join(join, " "))

	err := r.db.QueryRow(query, id).Scan(
		&conversation.ID,
		&conversation.TenantID,
		&conversation.CustomerID,
		&conversation.AssignedAgentID,
		&conversation.Status,
		&conversation.TenantName,
		&conversation.AgentName,
		&conversation.AgentEmail,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return conversation, nil
}

// GetByTenantAndCustomer implements [repository.ConversationRepository].
func (r *ConversationMySqlRepository) GetByTenantAndCustomer(tenantID uint, customerID uuid.UUID) (*domain.ConversationEntity, error) {
	var conversation domain.ConversationEntity
	err := r.db.QueryRow("SELECT id, tenant_id, customer_id, assigned_agent_id, status FROM conversations WHERE tenant_id = ? AND customer_id = UUID_TO_BIN(?)", tenantID, customerID.String()).Scan(&conversation.ID, &conversation.TenantID, &conversation.CustomerID, &conversation.AssignedAgentID, &conversation.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &conversation, nil
}

// UpdateStatus implements [repository.ConversationRepository].
func (r *ConversationMySqlRepository) UpdateStatus(id uuid.UUID, status string) error {
	query := "UPDATE conversations SET status = ? WHERE id = UUID_TO_BIN(?)"
	_, err := r.db.ExecContext(context.Background(), query, status, id)
	return err
}

// UpdateAssignment implements [repository.ConversationRepository].
func (r *ConversationMySqlRepository) UpdateAssignment(conv *domain.ConversationEntity, agentID uuid.UUID) error {
	query := "UPDATE conversations SET assigned_agent_id = UUID_TO_BIN(?) WHERE id = UUID_TO_BIN(?)"
	_, err := r.db.ExecContext(context.Background(), query, agentID, conv.ID)
	log.Printf("AgentID: %s\n", agentID.String())
	return err
}
