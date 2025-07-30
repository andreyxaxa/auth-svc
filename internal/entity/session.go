package entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	RefreshHash string    `json:"refresh_hash"`
	UserAgent   string    `json:"user_agent"`
	IP          string    `json:"ip"`
	CreatedAt   time.Time `json:"created_at"`
	Used        bool      `json:"used"`
}
