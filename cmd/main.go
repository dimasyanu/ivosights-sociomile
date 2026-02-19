package main

import (
	"github.com/dimasyanu/ivosights-sociomile/config"
	"github.com/dimasyanu/ivosights-sociomile/internal/delivery/rest"
)

func main() {
	restConfig := config.NewRestConfig()
	api := rest.NewRestApi(restConfig)
	api.Start()
}
