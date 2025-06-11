package app

import (
	"backend/internal/adapters/config"
	"backend/internal/adapters/controller/api/validator"
	"backend/internal/adapters/logger"
	"gorm.io/gorm"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v4"
)

// App is a struct that contains the fiber app, database connection, listen port, validator, logging boolean etc.
type App struct {
	Server    *echo.Echo
	Router    huma.API
	DB        *gorm.DB
	Validator *validator.Validator
}

// New is a function that creates a new app struct
func New(config *config.Config) *App {
	cfg := huma.DefaultConfig("backend", "v1.0.0")
	cfg.SchemasPath = "/docs#/schemas"
	cfg.OpenAPI.Servers = []*huma.Server{
		{
			URL:         "http://localhost:8080",
			Description: "local dev server",
		},
		{
			URL:         "https://back.lxft.tech",
			Description: "PROD",
		},
	}
	router := echo.New()
	api := humaecho.New(router, cfg)

	return &App{
		Server:    router,
		Router:    api,
		DB:        config.Database,
		Validator: validator.New(),
	}
}

// Start is a function that starts the app
func (a *App) Start() {
	if err := a.Server.Start(":8080"); err != nil {
		logger.Log.Panicf("failed to start listen (no tls): %v", err)
	}
}
