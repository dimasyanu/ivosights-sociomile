package mysqlrepo

import (
	"database/sql"

	"github.com/dimasyanu/ivosights-sociomile/internal/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/domain/repo"
)

type TenantMysqlRepository struct {
	db *sql.DB
}

func NewTenantRepository(db *sql.DB) repo.TenantRepository {
	return &TenantMysqlRepository{db: db}
}

func (r *TenantMysqlRepository) GetTenants(filter *domain.TenantFilter) ([]domain.TenantEntity, uint64, error) {
	var tenants []domain.TenantEntity

	// Build the base query
	query := "SELECT id, name, created_at, updated_at FROM tenants WHERE deleted_at IS NULL"
	countQuery := "SELECT COUNT(*) FROM tenants WHERE deleted_at IS NULL"

	// Add filtering conditions if needed (e.g., by name)
	if filter.Name != "" {
		query += " AND name LIKE ?"
		countQuery += " AND name LIKE ?"
	}

	// Add pagination
	query += " LIMIT ? OFFSET ?"

	// Prepare arguments for the queries
	var args []interface{}
	if filter.Name != "" {
		args = append(args, "%"+filter.Name+"%")
	}
	args = append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)

	// Execute count query to get total number of records
	var total uint64
	err := r.db.QueryRow(countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Execute main query to get paginated results
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var tenant domain.TenantEntity
		err := rows.Scan(&tenant.ID, &tenant.Name, &tenant.CreatedAt, &tenant.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		tenants = append(tenants, tenant)
	}

	return tenants, total, nil
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
