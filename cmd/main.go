package main

import (
	"log"

	"github.com/KarmaBeLike/time-tracker-api/config"
	postgres "github.com/KarmaBeLike/time-tracker-api/internal/database"
	"github.com/KarmaBeLike/time-tracker-api/internal/handlers"
	repositories "github.com/KarmaBeLike/time-tracker-api/internal/repository"
	"github.com/KarmaBeLike/time-tracker-api/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := postgres.OpenDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Выполнение миграций
	if err := postgres.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	router := gin.Default()

	// Инициализация репозитория и сервиса
	userRepo := repositories.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Регистрация маршрутов
	userHandler.Routes(router, cfg)

	// Проверка доступности главной страницы
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not Found"})
	})

	log.Println("Server is running at http://localhost:8080")
	router.Run(":8080")
}
