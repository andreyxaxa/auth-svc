package usecase

import (
	"context"

	"github.com/andreyxaxa/auth-svc/internal/entity"
	"github.com/google/uuid"
)

type (
	Sessions interface {
		Create(context.Context, uuid.UUID, string, string) (entity.Token, error)
		Refresh(context.Context, uuid.UUID, string, string, string) (entity.Token, error)
		Logout(context.Context, uuid.UUID) error
		Me(context.Context, string) (uuid.UUID, error)
	}
)
