package models

import "time"

type Reminder struct {
	ID        int64         `db:"id"`
	Title     string        `db:"title"`
	StartTime time.Time     `db:"start_time"`
	EndTime   time.Time     `db:"end_time"`
	Frequency time.Duration `db:"frequency"`
	ChatID    int64         `db:"chat_id"`
}
