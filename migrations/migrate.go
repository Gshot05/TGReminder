package migrations

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Функция написана только ради тестов чтобы каждый раз не дрочить терминал
// К моменту появления продовой версии она будет вырезана и будут настроены
// Человеческие автомиграции

func RunInitSQL() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	content, err := os.ReadFile("migrations/init.sql")
	if err != nil {
		log.Fatalf("failed to read init.sql: %v", err)
	}

	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer dbpool.Close()

	if _, err := dbpool.Exec(context.Background(), string(content)); err != nil {
		log.Fatalf("failed to execute init.sql: %v", err)
	}

	fmt.Println("✅ init.sql executed successfully")
}
