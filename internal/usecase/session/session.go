package session

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/andreyxaxa/auth-svc/internal/entity"
	"github.com/andreyxaxa/auth-svc/internal/repo"
	tokenmn "github.com/andreyxaxa/auth-svc/internal/token"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrSessionUsed       = errors.New("refresh token already used")
	ErrUserAgentMismatch = errors.New("user-agent mismath")
)

type UseCase struct {
	repo repo.SessionRepo
	tm   tokenmn.TokenManager

	whURL string
}

func New(r repo.SessionRepo, t tokenmn.TokenManager, wh string) *UseCase {
	return &UseCase{
		repo:  r,
		tm:    t,
		whURL: wh,
	}
}

func (uc *UseCase) Create(ctx context.Context, userID uuid.UUID, ua string, ip string) (entity.Token, error) {
	accessToken, err := uc.tm.Generate(userID)
	if err != nil {
		return entity.Token{}, fmt.Errorf("UseCase - Create - uc.tm.Generate: %w", err)
	}

	refreshRaw, err := generateBase64(32)
	if err != nil {
		return entity.Token{}, fmt.Errorf("UseCase - Create - generateBase64: %w", err)
	}

	refreshHash, err := bcrypt.GenerateFromPassword([]byte(refreshRaw), bcrypt.DefaultCost)
	if err != nil {
		return entity.Token{}, fmt.Errorf("UseCase - Create - bcrypt.GenerateFromPassword: %w", err)
	}

	sess := &entity.Session{
		ID:          uuid.New(),
		UserID:      userID,
		RefreshHash: string(refreshHash),
		UserAgent:   ua,
		IP:          ip,
		CreatedAt:   time.Now(),
		Used:        false,
	}

	if err := uc.repo.Create(ctx, sess); err != nil {
		return entity.Token{}, fmt.Errorf("UseCase - Create - uc.repo.Create: %w", err)
	}

	return entity.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshRaw,
	}, nil
}

func (uc *UseCase) Refresh(ctx context.Context, userID uuid.UUID, refresh string, ua string, ip string) (entity.Token, error) {
	session, err := uc.repo.GetConcreteByUserIDAndRawToken(ctx, userID, refresh)
	if err != nil {
		return entity.Token{}, fmt.Errorf("UseCase - Refresh - uc.repo.GetConcreteByUserIDAndRawToken: %w", err)
	}

	if session.Used {
		err = uc.repo.DeleteAllByUserID(ctx, userID)
		if err != nil {
			return entity.Token{}, fmt.Errorf("UseCase - Refresh - uc.repo.DeleteAllByUserID: %w", err)
		}

		return entity.Token{}, ErrSessionUsed
	}

	if session.UserAgent != ua {
		err = uc.repo.DeleteAllByUserID(ctx, userID)
		if err != nil {
			return entity.Token{}, fmt.Errorf("UseCase - Refresh - uc.repo.DeleteAllByUserID: %w", err)
		}

		return entity.Token{}, ErrUserAgentMismatch
	}

	if session.IP != ip && uc.whURL != "" {
		go uc.sendWebhookNotification(userID, ip)
	}

	if err := uc.repo.MarkAsUsed(ctx, session.ID); err != nil {
		return entity.Token{}, fmt.Errorf("UseCase - Refresh - uc.repo.MarkAsUsed: %w", err)
	}

	newPair, err := uc.Create(ctx, userID, ua, ip)
	if err != nil {
		return entity.Token{}, fmt.Errorf("UseCase - Refresh - uc.Create: %w", err)
	}

	return newPair, nil
}

func (uc *UseCase) Logout(ctx context.Context, userID uuid.UUID) error {
	if err := uc.repo.DeleteAllByUserID(ctx, userID); err != nil {
		return fmt.Errorf("UseCase - Logout - uc.repo.DeleteAllByUserID: %w", err)
	}

	return nil
}

func (uc *UseCase) Me(ctx context.Context, token string) (uuid.UUID, error) {
	uid, err := uc.tm.Parse(token)
	if err != nil {
		return uuid.Nil, fmt.Errorf("UseCase - Me - uc.tm.Parse: %w", err)
	}

	sessions, err := uc.repo.GetAllByUserID(ctx, uid)
	if err != nil {
		return uuid.Nil, fmt.Errorf("UseCase - Me - uc.repo.GetAllByUserID: %w", err)
	}

	if len(sessions) == 0 {
		return uuid.Nil, errors.New("invalid session")
	}

	return uid, nil
}

func (uc *UseCase) sendWebhookNotification(userID uuid.UUID, ip string) {
	payload := map[string]string{
		"user_id": userID.String(),
		"new_ip":  ip,
		"time":    time.Now().Format(time.RFC3339),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, uc.whURL, bytes.NewBuffer(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	_, _ = client.Do(req)
}

func generateBase64(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}
