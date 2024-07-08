package main

import (
	"log"
	"os"

	"github.com/KarmaBeLike/time-tracker-api/config"
	postgres "github.com/KarmaBeLike/time-tracker-api/internal/database"
	"github.com/KarmaBeLike/time-tracker-api/internal/external"
	"github.com/KarmaBeLike/time-tracker-api/internal/handlers"
	repositories "github.com/KarmaBeLike/time-tracker-api/internal/repository"
	"github.com/KarmaBeLike/time-tracker-api/internal/service"
	"github.com/KarmaBeLike/time-tracker-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	logger.New(os.Stdout, logger.LevelDebug)

	cfg, err := config.Load()
	if err != nil {
		logger.PrintError(err, nil)
	}

	db, err := postgres.OpenDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()

	// Выполнение миграций
	if err := postgres.RunMigrations(db); err != nil {
		logger.PrintFatal(err, nil)
	}

	router := gin.Default()

	// Инициализация репозитория и сервиса
	userRepo := repositories.NewUserRepository(db)
	peopleAPIClient := external.NewPeopleAPIClient(cfg.PeopleAPIBaseURL)
	userService := service.NewUserService(userRepo, peopleAPIClient)
	userHandler := handlers.NewUserHandler(userService)

	taskRepo := repositories.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handlers.NewTaskHandler(taskService)

	// Регистрация маршрутов
	userHandler.Routes(router, cfg)
	taskHandler.Routes(router, cfg)

	// Проверка доступности главной страницы
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not Found"})
	})

	log.Println("Server is running at http://localhost:8080")
	router.Run(":8080")
}
