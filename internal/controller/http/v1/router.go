package v1

import (
	"github.com/andreyxaxa/auth-svc/internal/controller/http/middleware"
	"github.com/andreyxaxa/auth-svc/internal/usecase"
	"github.com/andreyxaxa/auth-svc/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func NewSessionRoutes(apiV1Group fiber.Router, s usecase.Sessions, l logger.Interface) {
	r := &V1{
		s: s,
		l: l,
	}

	sessionGroup := apiV1Group.Group("/session")
	userGroup := apiV1Group.Group("/user")

	{ // /session
		// open
		sessionGroup.Post("/token", r.create)
		sessionGroup.Post("/refresh", r.refresh)
		// safe
		protectedSession := sessionGroup.Use(middleware.Auth(s))

		protectedSession.Post("/logout", r.logout)
	}

	{ // /user
		// safe
		protectedUser := userGroup.Use(middleware.Auth(s))

		protectedUser.Get("/me", r.me)
	}
}
