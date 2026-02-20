package rest

import (
	"database/sql"

	"github.com/dimasyanu/ivosights-sociomile/config"
	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest/models"
	"github.com/dimasyanu/ivosights-sociomile/internal/repository/mysqlrepo"
	"github.com/dimasyanu/ivosights-sociomile/service"
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(app *fiber.App, db *sql.DB) {
	// Initialize repositories
	userRepo := mysqlrepo.NewUserRepository(db)

	// Initialize services
	jwtService := service.NewJwtService(config.NewJwtConfig(config.EnvPath))
	authService := service.NewAuthService(userRepo, jwtService)

	// Initialize handlers with their respective services and repositories
	authHandler := NewAuthHandler(authService)

	// == Public Routes ==

	app.Get("/", func(ctx fiber.Ctx) error {
		return ctx.JSON(&models.Res[any]{
			Status:  fiber.StatusOK,
			Message: "Success",
		})
	})

	api := app.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)
}
