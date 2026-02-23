package service

import (
	"github.com/dimasyanu/ivosights-sociomile/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/repository"
)

type TenantService struct {
	repo repository.TenantRepository
}

func NewTenantService(repo repository.TenantRepository) *TenantService {
	return &TenantService{repo: repo}
}

func (s *TenantService) GetTenantByID(tenantID string) (*domain.Tenant, error) {
	tenantEntity, err := s.repo.GetTenantByID(tenantID)
	if err != nil {
		return nil, err
	}

	return &domain.Tenant{
		ID:        tenantEntity.ID,
		Name:      tenantEntity.Name,
		CreatedAt: tenantEntity.CreatedAt,
		UpdatedAt: tenantEntity.UpdatedAt,
	}, nil
}

func (s *TenantService) CreateTenant(name string) (*domain.Tenant, error) {
	data := &domain.TenantEntity{
		Name: name,
	}
	entity, err := s.repo.CreateTenant(data)
	if err != nil {
		return nil, err
	}
	return entity.ToDto(), nil
}

func (s *TenantService) UpdateTenant(tenantID string, name string) (*domain.Tenant, error) {
	// Implementation for updating a tenant would go here
	return nil, nil
}

func (s *TenantService) DeleteTenant(tenantID string) error {
	// Implementation for deleting a tenant would go here
	return nil
}
