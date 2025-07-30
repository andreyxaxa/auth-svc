package tokenmn

import "github.com/google/uuid"

type TokenManager interface {
	Generate(uuid.UUID) (string, error)
	Parse(string) (uuid.UUID, error)
}
