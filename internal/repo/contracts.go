package repo

import (
	"context"

	"github.com/andreyxaxa/auth-svc/internal/entity"
	"github.com/google/uuid"
)

type SessionRepo interface {
	Create(context.Context, *entity.Session) error
	GetAllByUserID(context.Context, uuid.UUID) ([]*entity.Session, error)
	GetConcreteByUserIDAndRawToken(context.Context, uuid.UUID, string) (*entity.Session, error)
	DeleteAllByUserID(context.Context, uuid.UUID) error
	DeleteConcreteByID(context.Context, uuid.UUID) error
	MarkAsUsed(context.Context, uuid.UUID) error
}
