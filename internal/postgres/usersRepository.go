package postgres

import (
	"context"
	"habits/internal/users"

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

func (r *userRepository) GetUserByGoogleID(ctx context.Context, googleID int64) (int, error) {
	query := `
		SELECT (id, name, username, telegram_id, telegram_chat_id) FROM users 
		WHERE google_id = $1 
	`

	var userID int
	err := r.db.QueryRow(ctx, query, googleID).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *userRepository) RegisterUserByGoogleID(ctx context.Context, googleID int64, name, username string) (int, error) {
	query := `
		INSERT INTO users (name, username, google_id) VALUES ($1, $2, $3)
	`

	var userID int
	err := r.db.QueryRow(ctx, query, name, username, googleID).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *users.User) error {
	query := `UPDATE users SET name = $1, username = $2, google_id = $3, telegram_id = $4, telegram_chat_id = $5, updated_at = NOW() WHERE id = $6`
	_, err := r.db.Exec(ctx, query, user.Name, user.Username, user.ID, user.GoogleID, user.TelegramID, user.TelegramChatID, user.ID)
	if err != nil {
		return err
	}
	return nil
}