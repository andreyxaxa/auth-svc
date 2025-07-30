package http

import (
	"github.com/andreyxaxa/auth-svc/config"
	_ "github.com/andreyxaxa/auth-svc/docs" // Swagger docs.
	v1 "github.com/andreyxaxa/auth-svc/internal/controller/http/v1"
	"github.com/andreyxaxa/auth-svc/internal/usecase"
	"github.com/andreyxaxa/auth-svc/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @title   Auth service API
// @version 1.0
// @host    localhost:8080
// BasePath /v1

// @securityDefinitions.apikey TokenAuth
// @in                         header
// @name 					   Authorization
func NewRouter(app *fiber.App, cfg *config.Config, s usecase.Sessions, l logger.Interface) {
	// Swagger
	if cfg.Swagger.Enabled {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	// Routers
	apiV1Group := app.Group("/v1")
	{
		v1.NewSessionRoutes(apiV1Group, s, l)
	}
}
