package rest

import (
	"database/sql"
	"fmt"

	"github.com/dimasyanu/ivosights-sociomile/config"
	"github.com/gofiber/fiber/v3"
	"github.com/yokeTH/gofiber-scalar/scalar/v3"
)

type RestApi struct {
	app  *fiber.App
	port uint16
}

func NewRestApi(c *config.RestConfig, db *sql.DB) *RestApi {
	app := fiber.New()

	app.Use(func(c fiber.Ctx) error {
		return c.Next()
	})

	RegisterRoutes(app, db)

	cfg := scalar.Config{
		RawSpecUrl: "/openapi.json",
	}

	app.Get("/docs/*", scalar.New(cfg))

	return &RestApi{
		app:  app,
		port: c.Port,
	}
}

func (api *RestApi) Start() {
	api.app.Listen(":" + fmt.Sprint(api.port))
}
