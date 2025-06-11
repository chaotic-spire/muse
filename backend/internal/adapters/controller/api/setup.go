package api

import (
	"backend/cmd/app"
	"backend/internal/domain/dto"
	"context"
	"github.com/danielgtaylor/huma/v2"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"net/http"
)

func Setup(app *app.App) {
	app.Server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	/*
		app.Server.Use(swagger.New(swagger.Config{
				BasePath: "/api/v1",
				FilePath: "./docs/swagger.json",
				Path:     "./docs",
				Title:    "Swagger API Docs",
			}))
	*/

	if viper.GetBool("settings.debug") {
		app.Server.Use(middleware.Logger())
	}

	// recover from panic
	app.Server.Use(middleware.Recover())

	// Provide a minimal config for startup check
	huma.Register(app.Router, huma.Operation{
		OperationID: "ping",
		Method:      http.MethodGet,
		Path:        "/api/ping",
		Summary:     "Pong!",
		Description: "Check if server has started and running",
		Tags:        []string{"ping"},
	}, func(ctx context.Context, input *struct{}) (*dto.PingOutput, error) {
		resp := &dto.PingOutput{}
		resp.Body.Status = "Pong!"
		return resp, nil
	})

	// middlewareHandler := middlewares.NewMiddlewareHandler(app)
	//
	// Setup user routes
	userHandler := NewUserHandler(app)
	userHandler.Setup(app.Router)
}
