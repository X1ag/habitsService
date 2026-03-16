package habits

import (
	"context"
)

type Repository interface {
	CreateHabit(ctx context.Context, habit *Habit) error
	GetHabits(ctx context.Context, userID int) ([]Habit, error)
	GetHabit(ctx context.Context, label string) (*Habit, error)
	UpdateHabit(ctx context.Context, habit *Habit) error
	DeleteHabit(ctx context.Context, label string, userID int) error
	GetOrderIndex(ctx context.Context, userID int) (int, error)
}