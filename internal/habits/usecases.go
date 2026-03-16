package habits

import (
	"context"
	"errors"
	"log/slog"
)

type HabitUsecase struct {
	repo Repository
}

func NewHabitUsecase(repo Repository) *HabitUsecase {
	return &HabitUsecase{
		repo: repo,
	}
}

func (u *HabitUsecase) CreateHabit(ctx context.Context, habit *Habit) error {
	if habit == nil {
		slog.Error("habit is nil")
		return nil
	}
	if habit.UserID == 0 {
		slog.Error("user id is 0", "user_id", habit.UserID)
		return errors.New("user id is required")

	}
	if habit.KeyName == "" {
		slog.Error("key name is required, but it's empty", "user_id", habit.UserID)
		return errors.New("key name is required")
	} 
	if habit.Label == "" {
		slog.Error("label is required, but it's empty", "user_id", habit.UserID)
		return errors.New("label is required")
	} 
	if habit.Type == "" {
		habit.Type = "boolean"
	}

	orderIndex, err := u.repo.GetOrderIndex(ctx, habit.UserID)
	if err != nil {
		slog.Error("error getting order index", "error", err, "user_id", habit.UserID)
		return err
	}

	habit.OrderIndex = orderIndex + 1
	habit.Active = true

	return u.repo.CreateHabit(ctx, habit)
}

func (u *HabitUsecase) GetHabit(ctx context.Context, habitLabel string) (*Habit, error) {
	return u.repo.GetHabit(ctx, habitLabel)
}

func (u *HabitUsecase) GetHabits(ctx context.Context, userID int) ([]Habit, error) {
	return u.repo.GetHabits(ctx, userID)
}

func (u *HabitUsecase) UpdateHabit(ctx context.Context, habit *Habit) error {
	return u.repo.UpdateHabit(ctx, habit)
}

func (u *HabitUsecase) DeleteHabit(ctx context.Context, label string, userID int) error {
	return u.repo.DeleteHabit(ctx, label, userID)
}