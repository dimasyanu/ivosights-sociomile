package repo

import "github.com/dimasyanu/ivosights-sociomile/internal/domain"

type TenantRepository interface {
	GetTenants(filter *domain.TenantFilter) ([]domain.TenantEntity, int64, error)
	GetByID(tenantID uint) (*domain.TenantEntity, error)
	Create(data *domain.TenantEntity) (*domain.TenantEntity, error)
	Update(tenantID uint, data *domain.TenantEntity) (*domain.TenantEntity, error)
}
