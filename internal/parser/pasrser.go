package parser

import (
	"fmt"
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
	err := utils.ParseEmptyText(text)
	if err != nil {
		return models.Reminder{}, err
	}

	// Парсим ключ-значения через utils
	values, err := utils.ParseKeyValuePairs(text)
	if err != nil {
		return models.Reminder{}, err
	}

	// Валидируем обязательные поля через utils
	requiredFields := []string{"Название", "Дата старта", "Дата конца", "Частота"}
	if err := utils.ValidateRequiredFields(values, requiredFields); err != nil {
		return models.Reminder{}, err
	}

	// Парсим даты через utils
	startTime, err := utils.ParseTimeInLocation(values["Дата старта"], config.TimeLayout, config.Location)
	if err != nil {
		return models.Reminder{}, fmt.Errorf("неверный формат даты старта: используйте %s", config.TimeLayout)
	}

	endTime, err := utils.ParseTimeInLocation(values["Дата конца"], config.TimeLayout, config.Location)
	if err != nil {
		return models.Reminder{}, fmt.Errorf("неверный формат даты конца: используйте %s", config.TimeLayout)
	}

	// Валидируем временной интервал через utils
	if err := utils.ValidateTimeInterval(startTime, endTime); err != nil {
		return models.Reminder{}, err
	}

	// Парсим частоту через utils
	frequency, err := utils.ParseFrequency(values["Частота"])
	if err != nil {
		return models.Reminder{}, err
	}

	// Создаем напоминание
	reminder := models.Reminder{
		Title:     strings.TrimSpace(values["Название"]),
		StartTime: startTime,
		EndTime:   endTime,
		Frequency: frequency,
		ChatID:    chatID,
	}

	return reminder, nil
}
