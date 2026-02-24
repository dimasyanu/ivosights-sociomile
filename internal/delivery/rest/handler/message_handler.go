package handler

import (
	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest/models"
	"github.com/dimasyanu/ivosights-sociomile/internal/domain/constant"
	"github.com/dimasyanu/ivosights-sociomile/internal/service"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
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
// @Router /channel/webhook [post]
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

// SendMessageToConversation godoc
// @Summary Reply to an existing conversation
// @Description Send a message to an existing conversation by its ID.
// @Tags Conversations
// @Accept json
// @Produce json
// @Param id path string true "Conversation ID"
// @Param message body models.SendMessageRequest true "Message payload"
// @Success 200 {object} any
// @Router /api/v1/backoffice/conversations/{id}/messages [post]
func (h *MessageHandler) SendMessageToConversation(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	convID, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&models.Res[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid conversation ID",
			Data:    err.Error(),
		})
	}

	var req models.SendMessageRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&models.Res[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request payload",
			Data:    err.Error(),
		})
	}

	msg, err := h.svc.CreateMessageInConversation(convID, constant.SenderTypeAgent, req.Message)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&models.Res[any]{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to send message",
			Data:    err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&models.Res[any]{
		Status:  fiber.StatusOK,
		Message: "Message sent successfully",
		Data:    msg.ID,
	})
}
