package handler

import (
	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest/models"
	"github.com/dimasyanu/ivosights-sociomile/internal/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/domain/repo"
	"github.com/dimasyanu/ivosights-sociomile/internal/service"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type ConversationHandler struct {
	svc *service.ConversationService
}

func NewConversationHandler(svc *service.ConversationService) *ConversationHandler {
	return &ConversationHandler{svc: svc}
}

// GetConversations godoc
// @Summary Get list of conversations
// @Description Get paginated list of conversations with optional filters
// @Tags Conversations
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param status query string false "Filter by status (open, closed, assigned)"
// @Success 200 {object} models.Res[*domain.Paginated[domain.Conversation]]
// @Router /conversations [get]
func (h *ConversationHandler) GetConversations(ctx fiber.Ctx) error {
	f := &domain.ConversationFilter{}

	// Bind query parameters to filter struct
	if err := ctx.Bind().Query(f); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&models.Res[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid query parameters",
		})
	}

	// Fetch conversations based on filter
	convs, err := h.svc.GetList(f)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&models.Res[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to fetch conversations",
		})
	}

	return ctx.JSON(&models.Res[*domain.Paginated[domain.Conversation]]{
		Status:  fiber.StatusOK,
		Message: "Retrieved successfully",
		Data:    convs,
	})
}

// GetConversationByID godoc
// @Summary Get conversation by ID
// @Description Get a single conversation by its ID
// @Tags Conversations
// @Accept json
// @Produce json
// @Param id path string true "Conversation ID"
// @Success 200 {object} models.Res[*domain.Conversation]
// @Router /conversations/{id} [get]
func (h *ConversationHandler) GetConversationByID(ctx fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&models.Res[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid conversation ID",
		})
	}

	conv, err := h.svc.GetByID(id)
	if err != nil {
		if err == repo.ErrNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(&models.Res[any]{
				Status:  fiber.StatusNotFound,
				Message: "Conversation not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(&models.Res[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to fetch conversation",
		})
	}
	if conv == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(&models.Res[any]{
			Status:  fiber.StatusNotFound,
			Message: "Conversation not found",
		})
	}

	return ctx.JSON(&models.Res[*domain.ConversationDetail]{
		Status:  fiber.StatusOK,
		Message: "Retrieved successfully",
		Data:    conv,
	})
}

// UpdateConversationStatus godoc
// @Summary Update conversation status
// @Description Update the status of a conversation (open, closed, assigned)
// @Tags Conversations
// @Accept json
// @Produce json
// @Param id path string true "Conversation ID"
// @Param body body struct{Status string `json:"status" validate:"required,oneof=open closed assigned"`} true "New status"
// @Success 200 {object} models.Res[any]
// @Router /conversations/{id}/status [put]
func (h *ConversationHandler) UpdateConversationStatus(ctx fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&models.Res[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid conversation ID",
		})
	}

	var req struct {
		Status string `json:"status" validate:"required,oneof=open closed assigned"`
	}
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&models.Res[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	err = h.svc.UpdateStatus(id, req.Status)
	if err != nil {
		if err == repo.ErrNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(&models.Res[any]{
				Status:  fiber.StatusNotFound,
				Message: "Conversation not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(&models.Res[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to update conversation status",
		})
	}

	return ctx.JSON(&models.Res[any]{
		Status:  fiber.StatusOK,
		Message: "Status updated successfully",
	})
}

// DeleteConversation godoc
// @Summary Delete conversation
// @Description Delete a conversation by its ID
// @Tags Conversations
// @Accept json
// @Produce json
// @Param id path string true "Conversation ID"
// @Success 200 {object} models.Res[any]
// @Router /conversations/{id} [delete]
func (h *ConversationHandler) DeleteConversation(ctx fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&models.Res[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid conversation ID",
		})
	}

	err = h.svc.Delete(id)
	if err != nil {
		if err == repo.ErrNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(&models.Res[any]{
				Status:  fiber.StatusNotFound,
				Message: "Conversation not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(&models.Res[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to delete conversation",
		})
	}

	return ctx.JSON(&models.Res[any]{
		Status:  fiber.StatusOK,
		Message: "Conversation deleted successfully",
	})
}
