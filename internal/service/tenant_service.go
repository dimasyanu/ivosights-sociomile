package service

import (
	"time"

	"github.com/dimasyanu/ivosights-sociomile/internal/domain"
	repository "github.com/dimasyanu/ivosights-sociomile/internal/domain/repo"
)

type TenantService struct {
	repo repository.TenantRepository
}

func NewTenantService(repo repository.TenantRepository) *TenantService {
	return &TenantService{repo: repo}
}

func (s *TenantService) GetTenants(filter *domain.TenantFilter) (*domain.Paginated[domain.Tenant], error) {
	tenantEntities, total, err := s.repo.GetTenants(filter)
	if err != nil {
		return nil, err
	}

	tenants := make([]domain.Tenant, len(tenantEntities))
	for i, entity := range tenantEntities {
		tenants[i] = *entity.ToDto()
	}

	return &domain.Paginated[domain.Tenant]{
		Items:    tenants,
		Page:     filter.Page,
		PageSize: filter.PageSize,
		Total:    total,
	}, nil
}

func (s *TenantService) GetTenantByID(id uint) (*domain.Tenant, error) {
	tenantEntity, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return tenantEntity.ToDto(), nil
}

func (s *TenantService) Create(name string) (*domain.Tenant, error) {
	data := &domain.TenantEntity{
		Name: name,
	}
	entity, err := s.repo.Create(data)
	if err != nil {
		return nil, err
	}
	return entity.ToDto(), nil
}

func (s *TenantService) Update(id uint, name string) (*domain.Tenant, error) {
	tenantEntity, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	tenantEntity.Name = name

	updatedEntity, err := s.repo.Update(id, tenantEntity)
	if err != nil {
		return nil, err
	}

	return updatedEntity.ToDto(), nil
}

func (s *TenantService) Delete(id uint) error {
	tenantEntity, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	now := time.Now()
	tenantEntity.DeletedAt = &now

	_, err = s.repo.Update(id, tenantEntity)
	if err != nil {
		return err
	}

	return nil
}
