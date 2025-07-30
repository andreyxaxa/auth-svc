package jwtmn

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken            = errors.New("invalid token")
	ErrInvalidClaims           = errors.New("invalid claims")
	ErrMissingSubClaim         = errors.New("missing 'sub' claim")
)

// JWT Token Manager
type TokenManager struct {
	secret   []byte
	duration time.Duration
}

func New(s string, d time.Duration) *TokenManager {
	return &TokenManager{
		secret:   []byte(s),
		duration: d,
	}
}

func (tm *TokenManager) Generate(userID uuid.UUID) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"iat": now.Unix(),
		"exp": now.Add(tm.duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signed, err := token.SignedString(tm.secret)
	if err != nil {
		return "", fmt.Errorf("TokenManager - Generate - token.SignedString: %w", err)
	}

	return signed, nil
}

func (tm *TokenManager) Parse(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok || t.Method.Alg() != jwt.SigningMethodHS512.Alg() {
			return nil, ErrUnexpectedSigningMethod
		}
		return tm.secret, nil
	})

	if err != nil || !token.Valid {
		return uuid.Nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, ErrInvalidClaims
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, ErrMissingSubClaim
	}

	uid, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, fmt.Errorf("TokenManager - Parse - uuid.Parse: %w", err)
	}

	return uid, nil
}
