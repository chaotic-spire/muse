package main

import (
	"backend/cmd/app"
	"backend/internal/adapters/config"
	"backend/internal/adapters/controller/api"
)

func main() {
	appConfig := config.Configure()
	mainApp := app.New(appConfig)

	api.Setup(mainApp)
	mainApp.Start()
}
