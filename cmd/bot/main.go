package main

import (
	"context"
	"log"
	"os"
	"tgreminder/internal/handler"
	"tgreminder/internal/repo"
	"tgreminder/internal/service"
	"tgreminder/migrations"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env файл не найден")
	}

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	// Удалить потом нахуй
	migrations.RunInitSQL()

	repo := repo.NewReminderRepo(dbpool)
	svc := service.NewReminderService(repo, bot)
	h := handler.NewHandler(svc)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		h.HandleUpdate(update)
	}
}
