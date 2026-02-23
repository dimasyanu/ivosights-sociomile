package handler

import (
	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest/models"
	"github.com/dimasyanu/ivosights-sociomile/internal/service"
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
// @Router       /api/v1/auth/login [post]
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
		Status:  fiber.StatusOK,
		Message: "Login successful",
		Data: models.LoginResponse{
			AccessToken: token,
		},
	})
}

func (h *AuthHandler) AuthorizationMiddleware(ctx fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&models.Res[any]{
			Status:  fiber.StatusUnauthorized,
			Message: "Missing Authorization header",
		})
	}

	token := authHeader[len("Bearer "):]
	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&models.Res[any]{
			Status:  fiber.StatusUnauthorized,
			Message: "Invalid Authorization header format",
		})
	}

	user, err := h.svc.ValidateToken(token)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&models.Res[any]{
			Status:  fiber.StatusUnauthorized,
			Message: "Invalid or expired token",
		})
	}

	ctx.Locals("user", user)
	return ctx.Next()
}
