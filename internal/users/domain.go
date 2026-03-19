package users

import "time"

type User struct {
	ID int `json:"id"`
	TelegramChatID int64 `json:"telegram_chat_id"`
	TelegramID int64 `json:"telegram_id"`
	GoogleID int64 `json:"google_id"`
	Email string `json:"email"`
	Name string `json:"name"`
	Username string `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Session struct {
	UserID int `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Claims struct {
		Sub string `json:"sub"`
		Email string `json:"email"`
		Name string `json:"name"`
		EmailVerified bool `json:"email_verified"`
		Picture string `json:"picture"`
	}