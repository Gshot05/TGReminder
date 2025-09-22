package service

import (
	"context"
	"fmt"
	"log"
	"tgreminder/internal/models"
	"tgreminder/internal/repo"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ReminderService struct {
	repo *repo.ReminderRepo
	bot  *tgbotapi.BotAPI
}

func NewReminderService(repo *repo.ReminderRepo, bot *tgbotapi.BotAPI) *ReminderService {
	return &ReminderService{repo: repo, bot: bot}
}

func (s *ReminderService) AddReminder(ctx context.Context, r models.Reminder) error {
	if err := s.repo.Save(ctx, r); err != nil {
		return err
	}

	go s.runReminder(r)
	return nil
}

func (s *ReminderService) runReminder(r models.Reminder) {
	now := time.Now()

	if now.After(r.EndTime) {
		return
	}

	if now.Before(r.StartTime) {
		waitTime := time.Until(r.StartTime)
		time.Sleep(waitTime)
	}

	if time.Now().Before(r.EndTime) {
		s.SendMessage(r.ChatID, fmt.Sprintf("⏰ Напоминание: %s", r.Title))
	}

	ticker := time.NewTicker(r.Frequency)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			currentTime := time.Now()
			if currentTime.After(r.EndTime) {
				if currentTime.Before(r.EndTime.Add(r.Frequency)) {
					s.SendMessage(r.ChatID, fmt.Sprintf("⏰ Последнее напоминание: %s", r.Title))
				}
				return
			}
			if currentTime.After(r.StartTime) {
				s.SendMessage(r.ChatID, fmt.Sprintf("⏰ Напоминание: %s", r.Title))
			}
		case <-time.After(time.Until(r.EndTime) + time.Minute):
			return
		}
	}
}

func (s *ReminderService) SendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := s.bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	} else {
		log.Printf("Message sent successfully")
	}
}
