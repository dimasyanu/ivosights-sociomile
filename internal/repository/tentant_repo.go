package repository

import "github.com/dimasyanu/ivosights-sociomile/domain"

type TenantRepository interface {
	GetTenantByID(tenantID string) (*domain.TenantEntity, error)
	CreateTenant(data *domain.TenantEntity) (*domain.TenantEntity, error)
	UpdateTenant(tenantID string, data *domain.TenantEntity) (*domain.TenantEntity, error)
}
