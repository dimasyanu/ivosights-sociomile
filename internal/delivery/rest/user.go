package rest

import (
	"fmt"

	"github.com/dimasyanu/ivosights-sociomile/config"
	"github.com/gofiber/fiber/v3"
)

type RestApi struct {
	app  *fiber.App
	port uint16
}

func NewRestApi(c *config.RestConfig) *RestApi {
	app := fiber.New()
	return &RestApi{
		app:  app,
		port: c.Port,
	}
}

func (api *RestApi) Start() {
	api.app.Listen(":" + fmt.Sprint(api.port))
}
