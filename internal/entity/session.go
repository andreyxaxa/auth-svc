package entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	RefreshHash string
	UserAgent   string
	IP          string
	CreatedAt   time.Time
	Used        bool
}
