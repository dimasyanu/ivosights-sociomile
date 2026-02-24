package mysqlrepo

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/dimasyanu/ivosights-sociomile/internal/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/domain/repo"
	"github.com/google/uuid"
)

type TicketMysqlRepo struct {
	db *sql.DB
}

func NewTicketMysqlRepo(db *sql.DB) repo.TicketRepository {
	return &TicketMysqlRepo{db: db}
}

// GetList implements [repo.TicketRepository].
func (r *TicketMysqlRepo) GetList(f *domain.TicketFilter) ([]*domain.TicketEntity, uint64, error) {
	panic("unimplemented")
}

// Create implements [repo.TicketRepository].
func (r *TicketMysqlRepo) Create(e *domain.TicketEntity) (*domain.TicketEntity, error) {
	id := uuid.New()
	pairs := map[string]any{
		"id":              id,
		"conversation_id": e.ConversationID,
		"title":           e.Title,
		"description":     e.Description,
		"status":          e.Status,
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
	e.ID = id
	return e, nil

}

func (r *TicketMysqlRepo) get(m map[string]any) (*domain.TicketEntity, error) {
	if len(m) == 0 {
		return nil, errors.New("no filter provided")
	}
	query := "SELECT id, conversation_id, title, description, status, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by FROM tickets WHERE "
	var args []any
	var conditions []string
	for k, v := range m {
		conditions = append(conditions, k+" = ?")
		args = append(args, v)
	}
	query += strings.Join(conditions, " AND ")

	row := r.db.QueryRowContext(context.Background(), query, args...)
	var t domain.TicketEntity
	if err := row.Scan(&t.ID, &t.ConversationID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.CreatedBy, &t.UpdatedAt, &t.UpdatedBy, &t.DeletedAt, &t.DeletedBy); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &t, nil
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
	query := "UPDATE tickets SET " + sets + " WHERE id = ?"
	vals = append(vals, e.ID)
	if _, err := r.db.ExecContext(context.Background(), query, vals...); err != nil {
		return nil, err
	}

	return e, nil
}
