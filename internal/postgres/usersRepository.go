package postgres

import (
	"context"
	"database/sql"
	"errors"
	"habits/internal/users"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUserByGoogle(ctx context.Context, googleID int64, name, username string) (int, error) {
	return 0, nil
}

func (r *userRepository) GetUser(ctx context.Context, telegramID int64) (*users.User, error) {
	return nil, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *users.User) (int, error) {
	query := `INSERT INTO users (telegram_id, name, username, email, google_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	
	var userID int
	err := r.db.QueryRow(ctx, query, user.TelegramID, user.Name, user.Username, user.GoogleID, user.Email).Scan(&userID)
	if err != nil {
		return -1, err
	}
	return userID, nil
}

func (r *userRepository) GetOrCreateUserID(ctx context.Context, telegramID int64, name, username string) (int, error) {
	query := `
		INSERT INTO users (telegram_id, name, username)
		VALUES ($1, $2, $3)
		ON CONFLICT (telegram_id)
		DO UPDATE SET
			name = $2,
			username = $3
		RETURNING id
	`

	var userID int
	err := r.db.QueryRow(ctx, query, telegramID, name, username).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// Auth Feature
func (r *userRepository) GetUserByGoogle(ctx context.Context, googleID int64) (int, error) {
	query := `
		SELECT (user_id, name, username, telegram_id, telegram_chat_id) FROM auth_identities 
		WHERE provider = 'google' AND provider_user_id = $1
	`

	var userID int
	err := r.db.QueryRow(ctx, query, googleID).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// AUTH Feature
func (r *userRepository) RegisterUserByGoogle(ctx context.Context, googleID int64, name, username string) (int, error) {
	query := `
		INSERT INTO auth_identities (user_id, provider, provider_user_id, email, username, display_name, avatar_url, last_login_at) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
	`

	var userID int
	err := r.db.QueryRow(ctx, query, name, username, googleID).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// AUTH Feature
func (r *userRepository) LoginWithGoogle(ctx context.Context, claims users.Claims, expiresAt time.Time) (*users.Session, error) {
	foundAuth := `
	SELECT user_id FROM auth_identities WHERE provider = 'google' AND provider_user_id = $1	
	`
	registerUser := `
		INSERT INTO users () VALUES () RETURNING id	
	`
	createSession := `
		INSERT INTO sessions (user_id, expires_at, last_seen_at) VALUES ($1, $2, NOW())	RETURNING id
	`
	createAuth := `
		INSERT INTO auth_identities (user_id, provider, provider_user_id, email, username, display_name, avatar_url, last_login_at) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())	
	`
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return nil, err
	}
	var googleID int
	err = tx.QueryRow(ctx, foundAuth, claims.Sub).Scan(&googleID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		var userID int
		err = tx.QueryRow(ctx, registerUser).Scan(&userID)
		if err != nil {
			tx.Rollback(ctx)
			return nil, err
		}
		_, err = tx.Exec(ctx, createAuth, userID, "google", claims.Sub, claims.Email, claims.Name, claims.Name, claims.Picture)
		if err != nil {
			tx.Rollback(ctx)
			return nil, err
		}
		_, err = tx.Exec(ctx, createSession, userID, expiresAt)
		if err != nil {
			tx.Rollback(ctx)
			return nil, err
		}
	case err != nil:
		tx.Rollback(ctx)
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	session := &users.Session{
		UserID:      googleID,
		ExpiresAt:   expiresAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return session, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *users.User) error {
	query := `UPDATE users SET name = $1, username = $2, google_id = $3, telegram_id = $4, telegram_chat_id = $5, updated_at = NOW() WHERE id = $6`
	_, err := r.db.Exec(ctx, query, user.Name, user.Username, user.ID, user.GoogleID, user.TelegramID, user.TelegramChatID, user.ID)
	if err != nil {
		return err
	}
	return nil
}