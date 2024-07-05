package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/KarmaBeLike/time-tracker-api/config"
	_ "github.com/lib/pq"

	"github.com/pkg/errors"
)

var DB *sql.DB

func OpenDB(cfg *config.Config) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, errors.Wrap(err, "open sql")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "connection is not established")
	}

	log.Println("Connected to DB")

	return db, nil
}

func RunMigrations(db *sql.DB) error {
	filePath := "migrations/create_users_table.up.sql"

	// Чтение SQL файла
	file, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read sql file: %w", err)
	}

	// Выполнение SQL команд
	_, err = db.Exec(string(file))
	if err != nil {
		return fmt.Errorf("execute sql file: %w", err)
	}

	log.Println("Migrations run successfully")
	return nil
}

// func executeSQLFile(db *sql.DB, filePath string) error {
// 	content, err := os.ReadFile(filePath)
// 	if err != nil {
// 		return errors.Wrap(err, "read sql file")
// 	}

// 	if _, err := db.Exec(string(content)); err != nil {
// 		return errors.Wrap(err, "execute sql file")
// 	}
// 	return nil
// }

func InitDB(cfg *config.Config) {
	var err error
	DB, err = OpenDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := RunMigrations(DB); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
}
