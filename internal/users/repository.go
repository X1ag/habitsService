package users

import (
	"context"
	"time"
)

type Repository interface {
	GetUser(ctx context.Context, telegramID int64) (*User, error)
	GetOrCreateUserID(ctx context.Context, telegramID int64, name, username string) (int, error)
	CreateUser(ctx context.Context, user *User) (int, error)
	GetUserByGoogle(ctx context.Context, googleID int64) (int, error)	
	UpdateUser(ctx context.Context, user *User) error
	LoginWithGoogle(ctx context.Context, claims Claims, expiresAt time.Time) (*Session, error) // auth

}