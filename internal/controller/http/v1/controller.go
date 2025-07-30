package v1

import (
	"github.com/andreyxaxa/auth-svc/internal/usecase"
	"github.com/andreyxaxa/auth-svc/pkg/logger"
)

type V1 struct {
	s usecase.Sessions
	l logger.Interface
}
