package persistent

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/andreyxaxa/auth-svc/internal/entity"
	"github.com/andreyxaxa/auth-svc/pkg/postgres"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

const (
	_defaultEntityCap = 64
	tableName         = "sessions"
)

type SessionRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *SessionRepo {
	return &SessionRepo{pg}
}

func (r *SessionRepo) Create(ctx context.Context, session *entity.Session) error {
	sql, args, err := r.Builder.
		Insert(tableName).
		Columns("id, user_id, refresh_hash, user_agent, ip, used").
		Values(session.ID, session.UserID, session.RefreshHash, session.UserAgent, session.IP, session.Used).
		ToSql()
	if err != nil {
		return fmt.Errorf("SessionRepo - Create - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("SessionRepo - Create - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *SessionRepo) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Session, error) {
	sql, args, err := r.Builder.
		Select("id, user_id, refresh_hash, user_agent, ip, created_at, used").
		From(tableName).
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("SessionRepo - GetByUserID - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("SessionRepo - GetByUserID - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	entities := make([]*entity.Session, 0, _defaultEntityCap)

	for rows.Next() {
		s := entity.Session{}

		err = rows.Scan(&s.ID, &s.UserID, &s.RefreshHash, &s.UserAgent, &s.IP, &s.CreatedAt, &s.Used)
		if err != nil {
			return nil, fmt.Errorf("SessionRepo - GetByUserID - rows.Scan: %w", err)
		}

		entities = append(entities, &s)
	}

	return entities, nil
}

func (r *SessionRepo) GetConcreteByUserIDAndRawToken(ctx context.Context, userID uuid.UUID, rawToken string) (*entity.Session, error) {
	sessions, err := r.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("SessionRepo - GetConcreteByUserIDAndRawToken - r.GetAllByUserID: %w", err)
	}

	for _, s := range sessions {
		if err := bcrypt.CompareHashAndPassword([]byte(s.RefreshHash), []byte(rawToken)); err == nil {
			return s, nil
		}
	}

	return nil, sql.ErrNoRows
}

func (r *SessionRepo) DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error {
	sql, args, err := r.Builder.
		Delete(tableName).
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("SessionRepo - DeleteAllByUserID - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("SessionRepo - DeleteAllByUserID - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *SessionRepo) DeleteConcreteByID(ctx context.Context, sessionID uuid.UUID) error {
	sql, args, err := r.Builder.
		Delete(tableName).
		Where(squirrel.Eq{"id": sessionID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("SessionRepo - DeleteConcreteByID - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("SessionRepo - DeleteConcreteByID - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *SessionRepo) MarkAsUsed(ctx context.Context, sessionID uuid.UUID) error {
	sql, args, err := r.Builder.
		Delete(tableName).
		Where(squirrel.Eq{"id": sessionID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("SessionRepo - MarkAsUsed - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("SessionRepo - MarkAsUsed - r.PoolExec: %w", err)
	}

	return nil
}
