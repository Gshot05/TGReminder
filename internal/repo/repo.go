package repo

import (
	"context"
	"tgreminder/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReminderRepo struct {
	db *pgxpool.Pool
}

func NewReminderRepo(db *pgxpool.Pool) *ReminderRepo {
	return &ReminderRepo{db: db}
}

func (r *ReminderRepo) Save(ctx context.Context, reminder models.Reminder) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO reminders (title, start_time, end_time, frequency, chat_id)
		VALUES ($1, $2, $3, $4, $5)`,
		reminder.Title,
		reminder.StartTime,
		reminder.EndTime,
		reminder.Frequency,
		reminder.ChatID,
	)
	return err
}
