package main

import (
	"log"
	"os"

	"github.com/KarmaBeLike/time-tracker-api/config"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	postgres "github.com/KarmaBeLike/time-tracker-api/internal/database"
	"github.com/KarmaBeLike/time-tracker-api/internal/external"
	"github.com/KarmaBeLike/time-tracker-api/internal/handlers"
	repositories "github.com/KarmaBeLike/time-tracker-api/internal/repository"
	"github.com/KarmaBeLike/time-tracker-api/internal/service"
	"github.com/KarmaBeLike/time-tracker-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

// @title Time Tracker API
// @version 1.0
// @description This is a sample time tracker server.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /

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
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.Run()

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
