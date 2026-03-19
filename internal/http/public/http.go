package http

import (
	"habits/internal/habits"
	"habits/internal/users"
)

type Handlers struct {
	userUsecase users.UsersUsecase
	habitUsecase habits.HabitUsecase
}

func NewUserHandlers(userUsecase users.UsersUsecase, habitUsecase habits.HabitUsecase) Handlers {
	return Handlers{
		userUsecase: userUsecase,
		habitUsecase: habitUsecase,
	}
}
