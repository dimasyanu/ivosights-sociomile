package handler

import (
	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest/models"
	"github.com/dimasyanu/ivosights-sociomile/internal/service"
	"github.com/gofiber/fiber/v3"
)

type MessageHandler struct {
	svc *service.MessageService
}

func NewMessageHandler(svc *service.MessageService) *MessageHandler {
	return &MessageHandler{svc: svc}
}

// HandleMessage godoc
// @Summary Handle incoming messages from channels
// @Description Process incoming messages from various channels.
// @Tags Messages
// @Accept json
// @Produce json
// @Param message body models.ChannelPayload true "Message payload"
// @Success 201 {object} any
// @Router /api/v1/message [post]
func (h *MessageHandler) HandleMessage(ctx fiber.Ctx) error {
	var req models.ChannelPayload
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&models.Res[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request payload",
			Data:    err.Error(),
		})
	}

	msg, err := h.svc.CreateMessage(req.TenantID, req.CustomerID, req.SenderType, req.Message)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&models.Res[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to process message",
			Data:    err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(&models.Res[any]{
		Status:  fiber.StatusCreated,
		Message: "Message processed successfully",
		Data:    msg.ID,
	})
}
