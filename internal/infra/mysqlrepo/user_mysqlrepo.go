package mysqlrepo

import (
	"context"
	"database/sql"
	"strings"

	"github.com/dimasyanu/ivosights-sociomile/internal/domain"
	repository "github.com/dimasyanu/ivosights-sociomile/internal/domain/repo"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type userMysqlRepository struct {
	db *sql.DB
}

var cols = []string{
	"id",
	"name",
	"email",
	"roles",
	"password_hash",
	"created_at",
	"created_by",
	"updated_at",
	"updated_by",
	"deleted_at",
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userMysqlRepository{
		db: db,
	}
}

func MapRowToUserEntity(row *sql.Row) (*domain.UserEntity, error) {
	var user domain.UserEntity
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Roles,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.CreatedBy,
		&user.UpdatedAt,
		&user.UpdatedBy,
		&user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userMysqlRepository) GetList(filter *domain.UserFilter) (*domain.Paginated[domain.UserEntity], int64, error) {

	// Get total count for pagination
	var count int64
	err := r.db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	// Build query with pagination
	query := "SELECT " + strings.Join(cols, ", ") + " FROM users"
	rows, err := r.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []domain.UserEntity
	for rows.Next() {
		var user domain.UserEntity
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Roles,
			&user.PasswordHash,
			&user.CreatedAt,
			&user.CreatedBy,
			&user.UpdatedAt,
			&user.UpdatedBy,
			&user.DeletedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	paginated := &domain.Paginated[domain.UserEntity]{
		Items:    users,
		Total:    count,
		Page:     filter.Page,
		PageSize: filter.PageSize,
	}
	return paginated, count, nil
}

func (r *userMysqlRepository) GetAvailableAgent() (*domain.UserEntity, error) {
	query := "SELECT " + strings.Join(cols, ", ") + " FROM users WHERE FIND_IN_SET('agent', roles) > 0 AND deleted_at IS NULL LIMIT 1"
	row := r.db.QueryRowContext(context.Background(), query)
	return MapRowToUserEntity(row)
}

func (r *userMysqlRepository) GetByID(id uuid.UUID) (*domain.UserEntity, error) {
	query := "SELECT " + strings.Join(cols, ", ") + " FROM users WHERE id = UUID_TO_BIN(?)"
	row := r.db.QueryRowContext(context.Background(), query, id.String())
	return MapRowToUserEntity(row)
}

func (r *userMysqlRepository) GetByEmail(email string) (*domain.UserEntity, error) {
	query := "SELECT " + strings.Join(cols, ", ") + " FROM users WHERE email = ?"
	row := r.db.QueryRowContext(context.Background(), query, email)
	return MapRowToUserEntity(row)
}

func (r *userMysqlRepository) Create(user *domain.UserEntity) (uuid.UUID, error) {
	id := uuid.New()
	pairs := map[string]any{
		"id":            id,
		"name":          user.Name,
		"email":         user.Email,
		"password_hash": user.PasswordHash,
		"roles":         user.Roles,

		"created_at": user.CreatedAt,
		"created_by": user.CreatedBy,
		"updated_at": user.UpdatedAt,
		"updated_by": user.UpdatedBy,
	}
	cols, slots, vals := MapForCreate(pairs)

	query := "INSERT INTO users (" + cols + ") VALUES (" + slots + ")"
	_, err := r.db.ExecContext(context.Background(), query, vals...)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *userMysqlRepository) Update(user *domain.UserEntity, data map[string]any) (*domain.UserEntity, error) {
	var name, email, roles string
	query := "SELECT name, email, roles FROM users WHERE id = UUID_TO_BIN(?)"
	row := r.db.QueryRowContext(context.Background(), query, user.ID.String())
	if err := row.Scan(&name, &email, &roles); err != nil {
		if err == sql.ErrNoRows {
			return nil, fiber.ErrNotFound
		}
		return nil, err
	}

	pairs := map[string]any{
		"deleted_at": user.DeletedAt,
	}

	if !user.UpdatedAt.IsZero() {
		pairs["updated_at"] = user.UpdatedAt
		pairs["updated_by"] = user.UpdatedBy
	}

	if val, ok := data["name"]; ok && user.Name != name {
		pairs["name"] = val
	}

	if val, ok := data["email"]; ok && user.Email != email {
		pairs["email"] = val
	}

	if val, ok := data["roles"]; ok && user.Roles != roles {
		pairs["roles"] = val
	}

	cols, vals := MapForUpdate(pairs)

	query = "UPDATE users SET " + cols + " WHERE id = UUID_TO_BIN(?)"
	vals = append(vals, user.ID.String())
	if _, err := r.db.ExecContext(context.Background(), query, vals...); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userMysqlRepository) Delete(id uuid.UUID) error {
	query := "DELETE FROM users WHERE id = UUID_TO_BIN(?)"
	_, err := r.db.ExecContext(context.Background(), query, id.String())
	if err != nil {
		return err
	}
	return nil
}
