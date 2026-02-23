package rest

import (
	"database/sql"

	"github.com/dimasyanu/ivosights-sociomile/config"
	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest/handler"
	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest/models"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra/mysqlrepo"
	"github.com/dimasyanu/ivosights-sociomile/internal/service"
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(app *fiber.App, db *sql.DB, mq infra.QueueClient, envPath string) {
	// Initialize repositories
	userRepo := mysqlrepo.NewUserRepository(db)
	convRepo := mysqlrepo.NewConversationRepository(db)
	msgRepo := mysqlrepo.NewMessageRepository(db)

	// Initialize services
	jwtService := service.NewJwtService(config.NewJwtConfig(envPath))
	authService := service.NewAuthService(userRepo, jwtService)
	userService := service.NewUserService(userRepo)
	convService := service.NewConversationService(convRepo, mq)
	messageService := service.NewMessageService(convService, msgRepo, mq)

	// Initialize handlers with their respective services and repositories
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	msgHandler := handler.NewMessageHandler(messageService)

	// == Public Routes ==

	app.Get("/", func(ctx fiber.Ctx) error {
		return ctx.JSON(&models.Res[any]{
			Status:  fiber.StatusOK,
			Message: "Success",
		})
	})

	api := app.Group("/api/v1")

	api.Post("/channel/webhook", msgHandler.HandleMessage)

	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)

	// == Protected Routes ==

	boff := api.Group("/backoffice")
	boff.Use(authHandler.AuthorizationMiddleware)
	boff.Get("/users", userHandler.GetUsers)
	boff.Get("/users/:id", userHandler.GetUserByID)
	boff.Post("/users", userHandler.CreateUser)
	boff.Put("/users/:id", userHandler.UpdateUser)
	boff.Delete("/users/:id", userHandler.DeleteUser)
}
