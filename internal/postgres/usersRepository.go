package postgres 

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) users.Repository {
	return &users.Repository{
		db: db,
	}
}
func (r *botRepository) CreateUser(ctx context.Context, user *users.User) (int, error) {
	query := `INSERT INTO users (telegram_id, name, username) VALUES ($1, $2, $3) RETURNING id`
	
	var userID int
	err := r.db.QueryRow(ctx, query, user.TelegramID, user.Name, user.Username).Scan(&userID)
	if err != nil {
		return -1, err
	}
	return userID, nil
}

func (r *botRepository) GetOrCreateUserID(ctx context.Context, telegramID int64, name, username string) (int, error) {
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