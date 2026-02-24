package handler

import (
	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest/models"
	"github.com/dimasyanu/ivosights-sociomile/internal/domain"
	"github.com/dimasyanu/ivosights-sociomile/internal/service"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type TicketHandler struct {
	svc *service.TicketService
}

func NewTicketHandler(svc *service.TicketService) *TicketHandler {
	return &TicketHandler{svc: svc}
}

func (h *TicketHandler) GetTickets(ctx fiber.Ctx) error {
	f := &domain.TicketFilter{}
	// Bind query parameters to filter struct
	if err := ctx.Bind().Query(f); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&models.Res[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid query parameters",
		})
	}

	res, err := h.svc.GetList(f)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&models.Res[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to fetch tickets",
		})
	}

	return ctx.JSON(&models.Res[*domain.Paginated[domain.Ticket]]{
		Status:  fiber.StatusOK,
		Message: "Retrieved successfully",
		Data:    res,
	})
}

func (h *TicketHandler) UpdateTicketStatus(ctx fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&models.Res[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid ticket ID",
		})
	}

	payload := &models.UpdateTicketStatusRequest{}
	if err := ctx.Bind().Body(payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&models.Res[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request payload",
		})
	}

	res, err := h.svc.UpdateStatus(id, payload.Status)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&models.Res[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to update ticket status",
		})
	}

	return ctx.JSON(&models.Res[*domain.Ticket]{
		Status:  fiber.StatusOK,
		Message: "Status updated successfully",
		Data:    res,
	})
}
