package app

import (
	"backend/internal/adapters/controller/api/validator"
	"backend/internal/domain/utils"
	"context"
	"errors"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v4"
)

// App is a struct that contains the fiber app, database connection, listen port, validator, logging boolean etc.
type App struct {
	Server    *echo.Echo
	Router    huma.API
	DB        *pgxpool.Pool
	Validator *validator.Validator
	Logger    *zap.Logger
	JwtSecret string
}

type config struct {
	// DbUrl - Postgres Database connection string
	// Example - "postgres://username:password@localhost:5432/database_name"
	DbUrl     string
	JwtSecret string `env:"JWT_SECRET"`
	BotToken  string `env:"BOT_TOKEN"`

	DbHost     string `env:"POSTGRES_HOST" env-default:"localhost"`
	DbPort     string `env:"POSTGRES_PORT" env-default:"5432"`
	DbPassword string `env:"POSTGRES_PASSWORD" env-default:"password"`
	DbUser     string `env:"POSTGRES_USER" env-default:"user"`
	DbName     string `env:"POSTGRES_DB" env-default:"db"`
}

// New is a function that creates a new app struct
func New(logger *zap.Logger) (*App, error) {
	var cfg config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	cfg.DbUrl = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName)
	logger.Info("connecting to db: " + cfg.DbUrl)

	if cfg.JwtSecret == "" {
		return nil, errors.New("JwtSecret is REQUIRED not to be null")
	}

	if cfg.BotToken == "" {
		return nil, errors.New("BotToken is REQUIRED not to be null")
	}

	apiCfg := huma.DefaultConfig("backend", "v1.0.0")
	apiCfg.SchemasPath = "/docs#/schemas"
	apiCfg.OpenAPI.Servers = []*huma.Server{
		{
			URL:         "http://localhost:8080",
			Description: "local dev server",
		},
		// TODO: enable in PROD
		/*
			{
				URL:         "https://back.lxft.tech",
				Description: "PROD",
			},
		*/
	}
	router := echo.New()
	router.HideBanner = true
	router.HidePort = true
	api := humaecho.New(router, apiCfg)

	conn, err := utils.NewConnection(context.Background(), cfg.DbUrl)
	if err != nil {
		return nil, err
	}

	requestValidator := validator.New(cfg.BotToken)

	return &App{
		Server:    router,
		Router:    api,
		DB:        conn,
		Validator: requestValidator,
		Logger:    logger,
	}, nil
}

// Start is a function that starts the app
func (a *App) Start() error {
	a.Logger.Info("starting server on :8080")
	if err := a.Server.Start(":8080"); err != nil {
		return err
	}
	return nil
}

func (a *App) Shutdown() error {
	a.Logger.Info("stopping server")
	err := a.Server.Close()
	if err != nil {
		return err
	}

	a.DB.Close()

	return nil
}
