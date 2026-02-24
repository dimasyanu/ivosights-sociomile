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
	authHandler,
		userHandler,
		convHandler,
		msgHandler,
		tenantHandler := getHandlers(db, mq, envPath)

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
	boff.Patch("/users/:id", userHandler.UpdateUser)
	boff.Delete("/users/:id", userHandler.DeleteUser)

	boff.Get("/tenants", tenantHandler.GetTenants)
	boff.Post("/tenants", tenantHandler.CreateTenant)
	boff.Patch("/tenants/:id", tenantHandler.UpdateTenant)
	boff.Delete("/tenants/:id", tenantHandler.DeleteTenant)

	boff.Get("/conversations", convHandler.GetConversations)
	boff.Get("/conversations/:id", convHandler.GetConversationByID)
	boff.Patch("/conversations/:id/status", convHandler.UpdateConversationStatus)
	boff.Delete("/conversations/:id", convHandler.DeleteConversation)
}

// Initialize handlers
func getHandlers(
	db *sql.DB,
	mq infra.QueueClient,
	envPath string,
) (
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	convHandler *handler.ConversationHandler,
	msgHandler *handler.MessageHandler,
	tenantHandler *handler.TenantHandler,
) {
	// Initialize repositories
	userRepo := mysqlrepo.NewUserRepository(db)
	convRepo := mysqlrepo.NewConversationRepository(db)
	msgRepo := mysqlrepo.NewMessageRepository(db)
	tenantRepo := mysqlrepo.NewTenantRepository(db)

	// Initialize services
	jwtService := service.NewJwtService(config.NewJwtConfig(envPath))
	authService := service.NewAuthService(userRepo, jwtService)
	userService := service.NewUserService(userRepo)
	convService := service.NewConversationService(convRepo, mq)
	messageService := service.NewMessageService(convService, msgRepo, mq)
	tenantService := service.NewTenantService(tenantRepo)

	// Initialize handlers with their respective services and repositories
	authHandler = handler.NewAuthHandler(authService)
	userHandler = handler.NewUserHandler(userService)
	convHandler = handler.NewConversationHandler(convService)
	msgHandler = handler.NewMessageHandler(messageService)
	tenantHandler = handler.NewTenantHandler(tenantService)

	return
}
