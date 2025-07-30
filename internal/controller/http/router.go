package http

import (
	"github.com/andreyxaxa/auth-svc/config"
	v1 "github.com/andreyxaxa/auth-svc/internal/controller/http/v1"
	"github.com/andreyxaxa/auth-svc/internal/usecase"
	"github.com/andreyxaxa/auth-svc/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func NewRouter(app *fiber.App, cfg *config.Config, s usecase.Sessions, l logger.Interface) {
	// TODO: swagger

	// if cfg.Swagger.Enabled { ... }

	apiV1Group := app.Group("/v1")
	{
		v1.NewSessionRoutes(apiV1Group, s, l)
	}
}
