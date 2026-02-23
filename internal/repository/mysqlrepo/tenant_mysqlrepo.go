package mysqlrepo

import (
	"database/sql"

	"github.com/dimasyanu/ivosights-sociomile/domain"
)

type TenantMysqlRepository struct {
	db *sql.DB
}

func NewTenantRepository(db *sql.DB) *TenantMysqlRepository {
	return &TenantMysqlRepository{db: db}
}

func (r *TenantMysqlRepository) GetTenantByID(tenantID string) (*domain.TenantEntity, error) {
	// Implementation for fetching tenant by ID from MySQL would go here
	return nil, nil
}

func (r *TenantMysqlRepository) CreateTenant(data *domain.TenantEntity) (*domain.TenantEntity, error) {
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

func (r *TenantMysqlRepository) UpdateTenant(tenantID string, data *domain.TenantEntity) (*domain.TenantEntity, error) {
	// Implementation for updating a tenant in MySQL would go here
	return nil, nil
}
