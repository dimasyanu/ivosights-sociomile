package rest

import (
	"database/sql"
	"fmt"

	"github.com/dimasyanu/ivosights-sociomile/config"
	"github.com/dimasyanu/ivosights-sociomile/internal/infra"
	"github.com/gofiber/fiber/v3"
	"github.com/yokeTH/gofiber-scalar/scalar/v3"
)

type RestApi struct {
	App  *fiber.App
	port uint16
}

func SetupApp(db *sql.DB, mq infra.QueueClient) *fiber.App {
	app := fiber.New()

	RegisterRoutes(app, db, mq, config.EnvPath)

	// Add OpenAPI UI route
	cfg := scalar.Config{
		RawSpecUrl: "/openapi.json",
	}
	app.Get("/docs/*", scalar.New(cfg))

	return app
}

func NewRestApi(c *config.RestConfig, db *sql.DB, mq infra.QueueClient) *RestApi {
	return &RestApi{
		App:  SetupApp(db, mq),
		port: c.Port,
	}
}

func (api *RestApi) Start() {
	api.App.Listen(":" + fmt.Sprint(api.port))
}
