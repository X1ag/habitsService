package users

import "context"

type Repository interface {
	GetUser(ctx context.Context, telegramID int64) (*User, error)
	GetOrCreateUserID(ctx context.Context, telegramID int64, name, username string) (int, error)
	CreateUser(ctx context.Context, user *User) (int, error)
}