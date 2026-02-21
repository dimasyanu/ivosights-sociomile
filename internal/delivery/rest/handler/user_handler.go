package handler

import (
	"database/sql"
	"errors"

	"github.com/dimasyanu/ivosights-sociomile/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest/models"
	"github.com/dimasyanu/ivosights-sociomile/service"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// GetUsers godoc
// @Summary Get a list of users
// @Description Retrieve a paginated list of users with optional filtering by name and email.
// @Tags Users
// @Accept json
// @Produce json
// @Param search query string false "Search term for name or email"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Success 200 {object} models.Res[domain.Paginated[domain.User]]
// @Failure 400 {object} models.Res[any]
// @Failure 500 {object} models.Res[any]
// @Router /api/v1/backoffice/users [get]
func (h *UserHandler) GetUsers(ctx fiber.Ctx) error {
	f := &domain.UserFilter{}

	// Bind query parameters to filter struct
	if err := ctx.Bind().Query(f); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&models.Res[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid query parameters",
		})
	}

	// Call service to get paginated users
	users, err := h.svc.GetUsers(f)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&models.Res[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to retrieve users",
		})
	}

	// Return paginated users in response
	return ctx.JSON(&models.Res[any]{
		Status:  fiber.StatusOK,
		Message: "Users retrieved successfully",
		Data:    users,
	})
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieve a user by their unique ID.
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.Res[domain.User]
// @Failure 400 {object} models.Res[any]
// @Failure 500 {object} models.Res[any]
// @Router /api/v1/backoffice/users/{id} [get]
func (h *UserHandler) GetUserByID(ctx fiber.Ctx) error {
	// Extract user ID from path parameters
	idStr := ctx.Params("id")
	if idStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(&models.Res[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid user ID",
		})
	}

	// Parse user string ID to UUID
	id, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&models.Res[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid user ID",
		})
	}

	// Call service to get user by ID
	user, err := h.svc.GetUserByID(id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&models.Res[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to retrieve user",
		})
	}
	if user == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(&models.Res[any]{
			Status:  fiber.StatusNotFound,
			Message: "User not found",
		})
	}

	// Return user in response
	return ctx.JSON(&models.Res[any]{
		Status:  fiber.StatusOK,
		Message: "User retrieved successfully",
		Data:    user,
	})
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided details.
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.UserCreateRequest true "User creation request"
// @Success 200 {object} models.Res[domain.User]
// @Failure 400 {object} models.Res[any]
// @Failure 500 {object} models.Res[any]
// @Router /api/v1/backoffice/users [post]
func (h *UserHandler) CreateUser(ctx fiber.Ctx) error {
	// Get authenticated user
	pic := ctx.Locals("user").(*domain.User)
	if pic == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&models.Res[any]{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	// Bind request body to UserCreateRequest struct
	payload := new(models.UserCreateRequest)
	if err := ctx.Bind().Body(payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&models.Res[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request payload",
			Data:    err.Error(),
		})
	}

	// Validate that password and repeat password match
	if payload.Password != payload.RepeatPassword {
		return ctx.Status(fiber.StatusBadRequest).JSON(&models.Res[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Password and repeat password do not match",
		})
	}

	// Validate existing email
	existingUser, err := h.svc.GetUserByEmail(payload.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&models.Res[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to check existing email",
		})
	}
	if existingUser != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&models.Res[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Email already exists",
		})
	}

	// Call service to create user
	user, err := h.svc.CreateUser(
		payload.Name,
		payload.Email,
		payload.Password,
		payload.Roles,
		pic.Email,
	)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&models.Res[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to create user",
		})
	}

	// Return created user in response
	return ctx.Status(fiber.StatusCreated).JSON(&models.Res[any]{
		Status:  fiber.StatusCreated,
		Message: "User created successfully",
		Data:    user,
	})
}

// UpdateUser godoc
// @Summary Update an existing user
// @Description Update an existing user's details by their unique ID.
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.UserUpdateRequest true "User update request"
// @Success 200 {object} models.Res[domain.User]
// @Failure 400 {object} models.Res[any]
// @Failure 500 {object} models.Res[any]
// @Router /api/v1/backoffice/users/{id} [put]
func (h *UserHandler) UpdateUser(ctx fiber.Ctx) error {
	pic := ctx.Locals("user").(*domain.User)
	if pic == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&models.Res[any]{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	idStr := ctx.Params("id")
	if idStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(&models.Res[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid user ID",
		})
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&models.Res[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid user ID",
		})
	}

	payload := new(models.UserUpdateRequest)
	if err := ctx.Bind().Body(payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&models.Res[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request payload",
			Data:    err.Error(),
		})
	}

	user, err := h.svc.UpdateUser(
		id,
		payload.Name,
		payload.Email,
		payload.Roles,
		pic.Email,
	)
	if err != nil {
		if errors.Is(err, fiber.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(&models.Res[any]{
				Status:  fiber.StatusNotFound,
				Message: "User not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(&models.Res[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to update user",
		})
	}

	return ctx.JSON(&models.Res[any]{
		Status:  fiber.StatusOK,
		Message: "User updated successfully",
		Data:    user,
	})
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by their unique ID.
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.Res[any]
// @Failure 400 {object} models.Res[any]
// @Failure 500 {object} models.Res[any]
// @Router /api/v1/backoffice/users/{id} [delete]
func (h *UserHandler) DeleteUser(ctx fiber.Ctx) error {
	pic := ctx.Locals("user").(*domain.User)
	if pic == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&models.Res[any]{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	idStr := ctx.Params("id")
	if idStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(&models.Res[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid user ID",
		})
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&models.Res[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid user ID",
		})
	}

	if err := h.svc.DeleteUser(id, pic.Email); err != nil {
		if errors.Is(err, fiber.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(&models.Res[any]{
				Status:  fiber.StatusNotFound,
				Message: "User not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(&models.Res[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to delete user",
		})
	}

	return ctx.JSON(&models.Res[any]{
		Status:  fiber.StatusOK,
		Message: "User deleted successfully",
	})
}
