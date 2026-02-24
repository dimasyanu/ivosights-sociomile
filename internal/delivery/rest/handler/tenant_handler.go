package handler

import (
	"strconv"

	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest/models"
	"github.com/dimasyanu/ivosights-sociomile/internal/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/service"
	"github.com/gofiber/fiber/v3"
)

type TenantHandler struct {
	svc *service.TenantService
}

func NewTenantHandler(svc *service.TenantService) *TenantHandler {
	return &TenantHandler{svc: svc}
}

// GetTenants godoc
// @Summary Get a list of tenants
// @Description Retrieve a paginated list of tenants with optional filtering by name.
// @Tags Tenants
// @Accept json
// @Produce json
// @Param search query string false "Search term for tenant name"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Success 200 {object} domain.Tenant
// @Router /api/v1/backoffice/tenants [get]
func (h *TenantHandler) GetTenants(ctx fiber.Ctx) error {
	f := &domain.TenantFilter{}

	// Bind query parameters to filter struct
	if err := ctx.Bind().Query(f); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&models.Res[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid query parameters",
		})
	}

	// Call service to get paginated tenants
	tenants, err := h.svc.GetTenants(f)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&models.Res[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to retrieve tenants",
		})
	}

	// Return paginated tenants in response
	return ctx.JSON(&models.Res[any]{
		Status:  fiber.StatusOK,
		Message: "Tenants retrieved successfully",
		Data:    tenants,
	})
}

// CreateTenant godoc
// @Summary Create a new tenant
// @Description Create a new tenant with the specified name.
// @Tags Tenants
// @Accept json
// @Produce json
// @Param tenant body models.TenantCreateRequest true "Tenant creation payload"
// @Success 201 {object} domain.Tenant
// @Router /api/v1/backoffice/tenants [post]
func (h *TenantHandler) CreateTenant(ctx fiber.Ctx) error {
	req := &models.TenantCreateRequest{}
	if err := ctx.Bind().Body(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&models.Res[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request",
		})
	}

	tenant, err := h.svc.Create(req.Name)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&models.Res[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to create tenant",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(&models.Res[any]{
		Status:  fiber.StatusCreated,
		Message: "Tenant created successfully",
		Data:    tenant,
	})
}

// UpdateTenant godoc
// @Summary Update an existing tenant
// @Description Update the name of an existing tenant by its ID.
// @Tags Tenants
// @Accept json
// @Produce json
// @Param id path string true "Tenant ID"
// @Param tenant body models.TenantUpdateRequest true "Tenant update payload"
// @Success 200 {object} domain.Tenant
// @Router /api/v1/backoffice/tenants/{id} [patch]
func (h *TenantHandler) UpdateTenant(ctx fiber.Ctx) error {
	payload := &models.TenantUpdateRequest{}
	if err := ctx.Bind().Body(payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&models.Res[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request",
		})
	}

	id := ctx.Params("id")
	idInt, _ := strconv.Atoi(id)
	tenant, err := h.svc.Update(uint(idInt), payload.Name)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&models.Res[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to update tenant",
		})
	}

	return ctx.JSON(&models.Res[any]{
		Status:  fiber.StatusOK,
		Message: "Tenant updated successfully",
		Data:    tenant,
	})
}

// DeleteTenant godoc
// @Summary Delete a tenant
// @Description Delete an existing tenant by its ID.
// @Tags Tenants
// @Accept json
// @Produce json
// @Param id path string true "Tenant ID"
// @Success 200 {object} any
// @Router /api/v1/backoffice/tenants/{id} [delete]
func (h *TenantHandler) DeleteTenant(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	idInt, _ := strconv.Atoi(id)
	err := h.svc.Delete(uint(idInt))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&models.Res[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to delete tenant",
		})
	}

	return ctx.JSON(&models.Res[any]{
		Status:  fiber.StatusOK,
		Message: "Tenant deleted successfully",
	})
}
