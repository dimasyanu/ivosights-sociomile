package mysqlrepo

import (
	"context"
	"database/sql"
	"strings"

	"github.com/dimasyanu/ivosights-sociomile/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/repository"
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

func (r *userMysqlRepository) GetUsers(filter *domain.UserFilter) (*domain.Paginated[domain.UserEntity], int64, error) {

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

func (r *userMysqlRepository) GetUserByID(id uuid.UUID) (*domain.UserEntity, error) {
	query := "SELECT " + strings.Join(cols, ", ") + " FROM users WHERE id = UUID_TO_BIN(?)"
	row := r.db.QueryRowContext(context.Background(), query, id.String())
	return MapRowToUserEntity(row)
}

func (r *userMysqlRepository) GetUserByEmail(email string) (*domain.UserEntity, error) {
	query := "SELECT " + strings.Join(cols, ", ") + " FROM users WHERE email = ?"
	row := r.db.QueryRowContext(context.Background(), query, email)
	return MapRowToUserEntity(row)
}

func (r *userMysqlRepository) CreateUser(user *domain.UserEntity) (uuid.UUID, error) {
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

func (r *userMysqlRepository) UpdateUser(user *domain.UserEntity) error {
	pairs := map[string]any{
		"name":  user.Name,
		"email": user.Email,
		"roles": user.Roles,

		"created_at": user.CreatedAt,
		"created_by": user.CreatedBy,
		"updated_at": user.UpdatedAt,
		"updated_by": user.UpdatedBy,
	}
	cols, vals := MapForUpdate(pairs)

	query := "UPDATE users SET " + cols + " WHERE id = UUID_TO_BIN(?)"
	vals = append(vals, user.ID.String())
	_, err := r.db.ExecContext(context.Background(), query, vals...)
	if err != nil {
		return err
	}

	return nil
}

func (r *userMysqlRepository) DeleteUser(id uuid.UUID) error {
	query := "DELETE FROM users WHERE id = UUID_TO_BIN(?)"
	_, err := r.db.ExecContext(context.Background(), query, id.String())
	if err != nil {
		return err
	}
	return nil
}
