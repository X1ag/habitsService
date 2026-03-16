package habits

import "time"

type Habit struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`
	KeyName string `json:"key_name"`
	Type string `json:"type"`
	Label string `json:"label"`
	Emoji string `json:"emoji"`
	OrderIndex int `json:"order_index"`
	Active bool `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}