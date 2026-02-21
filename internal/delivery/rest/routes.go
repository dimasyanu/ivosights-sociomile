package rest

import (
	"database/sql"

	"github.com/dimasyanu/ivosights-sociomile/config"
	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest/handler"
	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest/models"
	"github.com/dimasyanu/ivosights-sociomile/internal/repository/mysqlrepo"
	"github.com/dimasyanu/ivosights-sociomile/service"
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(app *fiber.App, db *sql.DB, envPath string) {
	// Initialize repositories
	userRepo := mysqlrepo.NewUserRepository(db)

	// Initialize services
	jwtService := service.NewJwtService(config.NewJwtConfig(envPath))
	authService := service.NewAuthService(userRepo, jwtService)
	userService := service.NewUserService(userRepo)

	// Initialize handlers with their respective services and repositories
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

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

	boff := api.Group("/backoffice")
	boff.Use(authHandler.AuthorizationMiddleware)
	boff.Get("/users", userHandler.GetUsers)
	boff.Get("/users/:id", userHandler.GetUserByID)
	boff.Post("/users", userHandler.CreateUser)
	boff.Put("/users/:id", userHandler.UpdateUser)
	boff.Delete("/users/:id", userHandler.DeleteUser)

}
