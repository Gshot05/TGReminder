package handler

import (
	"context"
	"tgreminder/internal/parser"
	"tgreminder/internal/service"
	"tgreminder/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	svc *service.ReminderService
}

func NewHandler(svc *service.ReminderService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) HandleUpdate(update tgbotapi.Update) {
	// log.Printf("Received update: %+v", update)

	err := utils.CheckMessage(update)
	if err != nil {
		return
	}

	text := update.Message.Text
	// log.Printf("Processing text: %s", text)

	switch text {
	case "/start":
		// log.Println("Handling /start command")
		h.svc.SendMessage(update.Message.Chat.ID, "👋 Привет! Я бот-напоминалка.\n\n"+
			"Чтобы создать напоминание, отправь сообщение в формате:\n"+
			"Название: Сделать разминку\n"+
			"Дата старта: 2025-09-17 08:00\n"+
			"Дата конца: 2025-09-20 20:00\n"+
			"Частота: каждые 3 часа\n")
		return
	}

	// log.Println("Attempting to parse reminder")
	reminder, err := parser.ParseReminder(text, update.Message.Chat.ID)
	if err != nil {
		// log.Printf("Parse error: %v", err)
		h.svc.SendMessage(update.Message.Chat.ID, "❌ Ошибка: "+err.Error())
		return
	}

	// log.Printf("Parsed reminder: %+v", reminder)

	if err := h.svc.AddReminder(context.Background(), reminder); err != nil {
		// log.Printf("Add reminder error: %v", err)
		h.svc.SendMessage(update.Message.Chat.ID, "❌ Не удалось добавить напоминание: "+err.Error())
		return
	}

	// log.Println("Reminder added successfully")
	h.svc.SendMessage(update.Message.Chat.ID, "✅ Напоминание добавлено: "+reminder.Title)
}
