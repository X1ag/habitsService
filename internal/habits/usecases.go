package habits

import (
	"context"
	"errors"
	"log/slog"
	"strings"
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
	if strings.TrimSpace(habit.KeyName) == "" {
		slog.Error("key name is required, but it's empty", "user_id", habit.UserID)
		return errors.New("key name is required")
	} 
	if strings.TrimSpace(habit.Label) == "" {
		slog.Error("label is required, but it's empty", "user_id", habit.UserID)
		return errors.New("label is required")
	} 
	if strings.TrimSpace(habit.Type) == "" {
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

func (u *HabitUsecase) GetHabit(ctx context.Context, keyName string, userID int) (*Habit, error) {
	if keyName == "" {
		slog.Error("key name is required, but it's empty")
		return nil, errors.New("key name is required")
	}
	if userID == 0 {
		slog.Error("user id is required, but it's empty or zero")
		return nil, errors.New("user id is required")
	}
	return u.repo.GetHabit(ctx, keyName, userID)
}

func (u *HabitUsecase) GetHabits(ctx context.Context, userID int) ([]Habit, error) {
	if userID == 0 {
		slog.Error("user id is required, but it's empty or zero")
		return nil, errors.New("user id is required")
	}
	return u.repo.GetHabits(ctx, userID)
}

func (u *HabitUsecase) UpdateHabit(ctx context.Context, habit *Habit) error {
	if habit.UserID == 0 {
		slog.Error("user id is required, but it's empty or zero")
		return errors.New("user id is required")
	}
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
	currentHabit, err := u.repo.GetHabit(ctx, habit.KeyName, habit.UserID)
	if err != nil {
		slog.Error("error getting habit", "error", err, "user_id", habit.UserID, "key_name", habit.KeyName)
	}
	if habit.Label == "" {
		habit.Label = currentHabit.Label
	} 
	if habit.Type == "" {
		habit.Type = currentHabit.Type
	}
	if habit.Emoji == "" {
		habit.Emoji = currentHabit.Emoji
	}

	orderIndex, err := u.repo.GetOrderIndex(ctx, habit.UserID)
	if err != nil {
		slog.Error("error getting order index", "error", err, "user_id", habit.UserID)
		return err
	}
	habit.OrderIndex = orderIndex
	
	return u.repo.UpdateHabit(ctx, habit)
}

func (u *HabitUsecase) DeleteHabit(ctx context.Context, keyName string, userID int) error {
	if keyName == "" {
		slog.Error("key name is required, but it's empty")
		return errors.New("key name is required")
	}	
	if userID == 0 {
		slog.Error("user id is required, but it's empty or zero")
		return errors.New("user id is required")
	}
	return u.repo.DeleteHabit(ctx, keyName, userID)
}