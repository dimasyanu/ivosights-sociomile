package rest

import (
	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest/models"
	"github.com/dimasyanu/ivosights-sociomile/service"
	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        loginRequest body models.LoginRequest true "Login request"
// @Success      200 {object} models.LoginResponse
// @Failure      400 {object} any
// @Failure      401 {object} any
// @Router       /auth/login [post]
func (h *AuthHandler) Login(ctx fiber.Ctx) error {
	req := new(models.LoginRequest)
	if err := ctx.Bind().Body(req); err != nil || req.Email == "" || req.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(&models.Res[any]{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	token, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&models.Res[any]{
			Status:  fiber.StatusUnauthorized,
			Message: "Invalid email or password",
		})
	}

	return ctx.JSON(&models.Res[models.LoginResponse]{
		Status: fiber.StatusOK,
		Data: models.LoginResponse{
			Token: token,
		},
	})
}

// Logout godoc
// @Summary      User logout
// @Description  Invalidate user session (if using server-side sessions)
// @Tags         auth
// @Produce      json
// @Success      200 {object} any
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(ctx fiber.Ctx) error {
	return ctx.JSON(&models.Res[any]{
		Status:  fiber.StatusOK,
		Message: "Logged out successfully",
	})
}
