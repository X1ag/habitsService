package postgres

import (
	"context"
	"habits/internal/habits"

	"github.com/jackc/pgx/v5/pgxpool"
)

type habitRepository struct {
	db *pgxpool.Pool
}

func NewHabitsRepository(db *pgxpool.Pool) *habitRepository {
	return &habitRepository{
		db: db,
	}
}

func (r *habitRepository) CreateHabit(ctx context.Context, habit *habits.Habit) error {
	query := `INSERT INTO habits (user_id, key_name, type, label, emoji, order_index, active) 
						VALUES ($1, $2, $3, $4, $5, $6, $7)`
	
	_, err := r.db.Exec(ctx, query, habit.UserID, habit.KeyName, habit.Type, habit.Label, habit.Emoji, habit.OrderIndex, habit.Active)
	if err != nil {
		return err
	}
	return nil
}

func (r *habitRepository) GetHabit(ctx context.Context, habitLabel string) (*habits.Habit, error) {
	query := `SELECT id, user_id, key_name, type, label, emoji, order_index, active, created_at, updated_at 
						FROM habits WHERE label = $1`

	var habit habits.Habit
	err := r.db.QueryRow(ctx, query, habitLabel).Scan(&habit.ID, &habit.UserID, &habit.KeyName, &habit.Type, &habit.Label, &habit.Emoji, &habit.OrderIndex, &habit.Active, &habit.CreatedAt, &habit.UpdatedAt)
	if err != nil {
		return &habits.Habit{}, err
	}

	return &habit, nil
}

func (r *habitRepository) GetHabits(ctx context.Context, userID int) ([]habits.Habit, error) {
	query := `SELECT id, user_id, key_name, type, label, emoji, order_index, active, created_at, updated_at 
						FROM habits WHERE user_id = $1`
	
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allHabits []habits.Habit
	for rows.Next() {
		var singleHabit habits.Habit
		err := rows.Scan(&singleHabit.ID, &singleHabit.UserID, &singleHabit.KeyName, &singleHabit.Type, &singleHabit.Label, &singleHabit.Emoji, &singleHabit.OrderIndex, &singleHabit.Active, &singleHabit.CreatedAt, &singleHabit.UpdatedAt)
		if err != nil {
			return nil, err
		}
		allHabits = append(allHabits, singleHabit)
	}

	return allHabits, nil
}

func (r *habitRepository) UpdateHabit(ctx context.Context, habit *habits.Habit) error {
	query := `UPDATE habits SET key_name = $1, type = $2, label = $3, emoji = $4, order_index = $5, active = $6, updated_at = NOW() WHERE id = $7`
	_, err := r.db.Exec(ctx, query, habit.KeyName, habit.Type, habit.Label, habit.Emoji, habit.OrderIndex, habit.Active, habit.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *habitRepository) DeleteHabit(ctx context.Context, label string, userID int) error {
	query := `DELETE FROM habits WHERE label = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, label, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *habitRepository) GetOrderIndex(ctx context.Context, userID int) (int, error) {
	query := `SELECT COALESCE(MAX(order_index), -1) FROM habits WHERE user_id = $1`
	
	var orderIndex int
	err := r.db.QueryRow(ctx, query, userID).Scan(&orderIndex)
	if err != nil {
		return 0, err
	}
	
	return orderIndex, nil
}