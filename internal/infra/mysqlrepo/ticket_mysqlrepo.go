package mysqlrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dimasyanu/ivosights-sociomile/internal/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/domain/repo"
	"github.com/google/uuid"
)

type TicketMysqlRepo struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) repo.TicketRepository {
	return &TicketMysqlRepo{db: db}
}

// GetList implements [repo.TicketRepository].
func (r *TicketMysqlRepo) GetList(f *domain.TicketFilter) ([]*domain.TicketEntity, uint64, error) {
	panic("unimplemented")
}

// Create implements [repo.TicketRepository].
func (r *TicketMysqlRepo) Create(e *domain.TicketEntity) (*domain.TicketEntity, error) {
	pairs := map[string]any{
		"id":              e.ID,
		"tenant_id":       e.TenantID,
		"conversation_id": e.ConversationID,
		"title":           e.Title,
		"description":     e.Description,
		"status":          e.Status,
		"priority":        e.Priority,
		"created_at":      e.CreatedAt,
		"created_by":      e.CreatedBy,
		"updated_at":      e.UpdatedAt,
		"updated_by":      e.UpdatedBy,
	}
	cols, slots, vals := MapForCreate(pairs)

	query := "INSERT INTO tickets (" + cols + ") VALUES (" + slots + ")"
	_, err := r.db.ExecContext(context.Background(), query, vals...)
	if err != nil {
		return nil, err
	}
	return e, nil

}

func (r *TicketMysqlRepo) get(m map[string]any) (*domain.TicketEntity, error) {
	if len(m) == 0 {
		return nil, errors.New("no filter provided")
	}
	query := "SELECT id, tenant_id, conversation_id, title, description, status, priority, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by FROM tickets WHERE "
	var args []any
	var conditions []string
	for k, v := range m {
		args = append(args, v)
		if _, ok := v.(interface{ String() string }); ok && strings.Contains(fmt.Sprintf("%T", v), "UUID") {
			conditions = append(conditions, k+" = UUID_TO_BIN(?)")
			continue
		}
		conditions = append(conditions, k+" = ?")
	}
	query += strings.Join(conditions, " AND ")

	row := r.db.QueryRowContext(context.Background(), query, args...)
	t := &domain.TicketEntity{}
	if err := row.Scan(
		&t.ID,
		&t.TenantID,
		&t.ConversationID,
		&t.Title,
		&t.Description,
		&t.Status,
		&t.Priority,
		&t.CreatedAt,
		&t.CreatedBy,
		&t.UpdatedAt,
		&t.UpdatedBy,
		&t.DeletedAt,
		&t.DeletedBy,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repo.ErrNotFound
		}
		return nil, err
	}

	return t, nil
}

// GetByConversationID implements [repo.TicketRepository].
func (r *TicketMysqlRepo) GetByConversationID(convID uuid.UUID) (*domain.TicketEntity, error) {
	return r.get(map[string]any{
		"conversation_id": convID,
	})
}

// GetByID implements [repo.TicketRepository].
func (r *TicketMysqlRepo) GetByID(id uuid.UUID) (*domain.TicketEntity, error) {
	return r.get(map[string]any{
		"id": id,
	})
}

// Update implements [repo.TicketRepository].
func (r *TicketMysqlRepo) Update(e *domain.TicketEntity) (*domain.TicketEntity, error) {
	pairs := map[string]any{
		"title":       e.Title,
		"description": e.Description,
		"status":      e.Status,
		"updated_at":  e.UpdatedAt,
		"updated_by":  e.UpdatedBy,
	}
	sets, vals := MapForUpdate(pairs)
	query := "UPDATE tickets SET " + sets + " WHERE BIN_TO_UUID(id) = ?"
	vals = append(vals, e.ID.String())
	if _, err := r.db.ExecContext(context.Background(), query, vals...); err != nil {
		return nil, err
	}

	return e, nil
}

func (r *TicketMysqlRepo) UpdateStatus(id uuid.UUID, status string) (*domain.TicketEntity, error) {
	query := "UPDATE tickets SET status = ?, updated_at = ? WHERE BIN_TO_UUID(id) = ?"
	if _, err := r.db.ExecContext(context.Background(), query, status, time.Now(), id.String()); err != nil {
		return nil, err
	}
	return r.GetByID(id)
}
