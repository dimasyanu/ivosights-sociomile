package mysqlrepo

import (
	"database/sql"

	"github.com/dimasyanu/ivosights-sociomile/internal/domain"
)

type TenantMysqlRepository struct {
	db *sql.DB
}

func NewTenantRepository(db *sql.DB) *TenantMysqlRepository {
	return &TenantMysqlRepository{db: db}
}

func (r *TenantMysqlRepository) GetByID(id uint) (*domain.TenantEntity, error) {
	query := "SELECT id, name, created_at, updated_at FROM tenants WHERE id = ? AND deleted_at IS NULL"
	row := r.db.QueryRow(query, id)

	var tenant domain.TenantEntity
	err := row.Scan(&tenant.ID, &tenant.Name, &tenant.CreatedAt, &tenant.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No tenant found with the given ID
		}
		return nil, err
	}

	return &tenant, nil
}

func (r *TenantMysqlRepository) Create(data *domain.TenantEntity) (*domain.TenantEntity, error) {
	res, err := r.db.Exec(
		"INSERT INTO tenants (name) VALUES (?)",
		data.Name,
	)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	data.ID = uint(id)

	return data, nil
}

func (r *TenantMysqlRepository) Update(id uint, data *domain.TenantEntity) (*domain.TenantEntity, error) {
	query := "UPDATE tenants SET name = ?, updated_at = NOW() WHERE id = ? AND deleted_at IS NULL"
	_, err := r.db.Exec(query, data.Name, id)
	if err != nil {
		return nil, err
	}

	return data, nil
}
