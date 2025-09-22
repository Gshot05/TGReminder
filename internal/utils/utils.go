package utils

import (
	"fmt"
	"strconv"
	"strings"
	errors2 "tgreminder/internal/errors"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Проверяет что сообщение не пустое
func CheckMessage(update tgbotapi.Update) error {
	if update.Message == nil || update.Message.Text == "" {
		return errors2.ErrEmptyMessage
	}
	return nil
}

// ParseEmptyText парсит пустой текст
func ParseEmptyText(text string) error {
	if strings.TrimSpace(text) == "" {
		return errors2.ErrEmptyMessage
	}
	return nil
}

// ParseKeyValuePairs парсит текст в map ключ-значение
func ParseKeyValuePairs(text string) (map[string]string, error) {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	if len(lines) < 4 {
		return nil, errors2.ErrWrongFormat
	}

	values := make(map[string]string)

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("неверный формат в строке %d: ожидается 'Ключ: Значение'", i+1)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" {
			return nil, fmt.Errorf("пустой ключ в строке %d", i+1)
		}

		if _, exists := values[key]; exists {
			return nil, fmt.Errorf("дублирующийся ключ: %s", key)
		}

		values[key] = value
	}

	return values, nil
}

// ValidateRequiredFields проверяет наличие обязательных полей
func ValidateRequiredFields(values map[string]string, required []string) error {
	var missing []string

	for _, field := range required {
		if values[field] == "" {
			missing = append(missing, field)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("отсутствуют обязательные поля: %s", strings.Join(missing, ", "))
	}

	return nil
}

// ParseTimeInLocation парсит время в указанной локации
func ParseTimeInLocation(timeStr, layout string, loc *time.Location) (time.Time, error) {
	return time.ParseInLocation(layout, timeStr, loc)
}

// ValidateTimeInterval проверяет корректность временного интервала
func ValidateTimeInterval(start, end time.Time) error {
	if start.After(end) {
		return errors2.ErrWrongStartDate
	}

	if time.Until(start) < 0 {
		return errors2.ErrWrongStartDate2
	}

	return nil
}

// ParseFrequency парсит строку частоты в time.Duration
func ParseFrequency(freqStr string) (time.Duration, error) {
	freqStr = strings.ToLower(strings.TrimSpace(freqStr))

	freqStr = strings.TrimPrefix(freqStr, "каждые")
	freqStr = strings.TrimSpace(freqStr)

	if freqStr == "" {
		return 0, errors2.ErrEmptyFreq
	}

	parts := strings.Fields(freqStr)
	if len(parts) < 2 {
		return 0, errors2.ErrWrongFreq
	}

	number, err := strconv.Atoi(parts[0])
	if err != nil || number <= 0 {
		return 0, errors2.ErrWrongNumber
	}

	unit := normalizeTimeUnit(parts[1])
	duration, err := parseTimeUnit(unit, number)
	if err != nil {
		return 0, err
	}

	if duration < time.Minute {
		return 0, errors2.ErrToShortFreq
	}

	return duration, nil
}

// normalizeTimeUnit нормализует единицы измерения времени
func normalizeTimeUnit(unit string) string {
	unit = strings.ToLower(unit)

	switch {
	case strings.HasPrefix(unit, "мин"), strings.HasPrefix(unit, "м"):
		return "minute"
	case strings.HasPrefix(unit, "час"), strings.HasPrefix(unit, "ч"):
		return "hour"
	case strings.HasPrefix(unit, "сек"), strings.HasPrefix(unit, "с"):
		return "second"
	case strings.HasPrefix(unit, "дн"), strings.HasPrefix(unit, "д"):
		return "day"
	default:
		return unit
	}
}

// parseTimeUnit конвертирует единицу измерения в time.Duration
func parseTimeUnit(unit string, number int) (time.Duration, error) {
	switch unit {
	case "second", "sec", "s", "сек", "с":
		return time.Duration(number) * time.Second, nil
	case "minute", "min", "m", "мин", "м":
		return time.Duration(number) * time.Minute, nil
	case "hour", "hr", "h", "час", "ч":
		return time.Duration(number) * time.Hour, nil
	case "day", "d", "дн", "д":
		return time.Duration(number) * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("неизвестная единица времени: %s. Используйте: секунды, минуты, часы, дни", unit)
	}
}

// ParseAndValidateDates парсит и валидирует даты старта и окончания
func ParseAndValidateDates(startDateStr, endDateStr, layout string, loc *time.Location) (time.Time, time.Time, error) {
	startTime, err := ParseTimeInLocation(startDateStr, layout, loc)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("неверный формат даты старта: используйте %s", layout)
	}

	endTime, err := ParseTimeInLocation(endDateStr, layout, loc)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("неверный формат даты конца: используйте %s", layout)
	}

	if err := ValidateTimeInterval(startTime, endTime); err != nil {
		return time.Time{}, time.Time{}, err
	}

	return startTime, endTime, nil
}
