package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andreyxaxa/auth-svc/config"
	"github.com/andreyxaxa/auth-svc/internal/controller/http"
	"github.com/andreyxaxa/auth-svc/internal/repo/persistent"
	jwtmn "github.com/andreyxaxa/auth-svc/internal/token/jwt"
	"github.com/andreyxaxa/auth-svc/internal/usecase/session"
	"github.com/andreyxaxa/auth-svc/pkg/httpserver"
	"github.com/andreyxaxa/auth-svc/pkg/logger"
	"github.com/andreyxaxa/auth-svc/pkg/postgres"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Token Manager
	tm := jwtmn.New(cfg.JWT.SecretKey, time.Minute*time.Duration(cfg.JWT.TTL))

	// Use-Case
	sessionUseCase := session.New(
		persistent.New(pg),
		tm,
		cfg.Webhook.URL,
	)

	// HTTP Server
	httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port))
	http.NewRouter(httpServer.App, cfg, sessionUseCase, l)

	// Start Server
	httpServer.Start()

	// Waiting Signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
