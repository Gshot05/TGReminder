package errors2

import "errors"

var (
	ErrEmptyMessage    = errors.New("Сообщение не может быть пустым!")
	ErrWrongFormat     = errors.New("Некорректный формат: нужно минимум 4 строки")
	ErrWrongStartDate  = errors.New("Дата старта не может быть позже даты окончания")
	ErrWrongStartDate2 = errors.New("Дата старта не может быть в прошлом")
	ErrEmptyFreq       = errors.New("Строка частоты не может быть пустой")
	ErrWrongFreq       = errors.New("Неверный формат частоты. Пример: '5 минут' или 'каждые 2 часа'")
	ErrWrongNumber     = errors.New("Неверное число в частоте: должно быть положительным целым числом")
	ErrToShortFreq     = errors.New("Частота не может быть менее 1 минуты")
)
