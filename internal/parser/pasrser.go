package parser

import (
	"strings"
	"tgreminder/internal/models"
	"tgreminder/internal/utils"
	"time"
)

type Config struct {
	TimeLayout string
	Location   *time.Location
}

func DefaultConfig() Config {
	loc, _ := time.LoadLocation("Local")
	if loc == nil {
		loc = time.UTC
	}

	return Config{
		TimeLayout: "2006-01-02 15:04",
		Location:   loc,
	}
}

type ParseResult struct {
	Reminder models.Reminder
	Warnings []string
}

func ParseReminder(text string, chatID int64) (models.Reminder, error) {
	return ParseReminderWithConfig(text, chatID, DefaultConfig())
}

func ParseReminderWithConfig(text string, chatID int64, config Config) (models.Reminder, error) {
	if err := utils.ParseEmptyText(text); err != nil {
		return models.Reminder{}, err
	}

	values, err := utils.ParseKeyValuePairs(text)
	if err != nil {
		return models.Reminder{}, err
	}

	requiredFields := []string{"Название", "Дата старта", "Дата конца", "Частота"}
	if err := utils.ValidateRequiredFields(values, requiredFields); err != nil {
		return models.Reminder{}, err
	}

	startTime, endTime, err := utils.ParseAndValidateDates(
		values["Дата старта"],
		values["Дата конца"],
		config.TimeLayout,
		config.Location,
	)
	if err != nil {
		return models.Reminder{}, err
	}

	frequency, err := utils.ParseFrequency(values["Частота"])
	if err != nil {
		return models.Reminder{}, err
	}

	return models.Reminder{
		Title:     strings.TrimSpace(values["Название"]),
		StartTime: startTime,
		EndTime:   endTime,
		Frequency: frequency,
		ChatID:    chatID,
	}, nil
}
