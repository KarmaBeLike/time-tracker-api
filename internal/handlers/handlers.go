package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/KarmaBeLike/time-tracker-api/config"
	"github.com/KarmaBeLike/time-tracker-api/internal/models"
	"github.com/KarmaBeLike/time-tracker-api/internal/service"
	"github.com/KarmaBeLike/time-tracker-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Routes(router *gin.Engine, cfg *config.Config) {
	user := router.Group("/users")
	{
		user.POST("/", h.CreateUser)
		user.GET("/", h.GetUsers)
		user.DELETE("/:userId", h.DeleteUser)
		user.PUT("/:userId", h.UpdateUser)
		user.GET("/:userId/worklogs", h.GetWorklogs)
		user.POST("/:userId/tasks/:taskId/start", h.StartTask)
		user.POST("/:userId/tasks/:taskId/stop", h.StopTask)

	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	// var json struct {
	// 	Name           string `json:"name"`
	// 	Surname        string `json:"surname"`
	// 	Patronymic     string `json:"patronymic"`
	// 	PassportNumber string `json:"passportNumber"`
	// 	Address        string `json:"address"`
	// }
	user := &models.User{}

	if err := c.ShouldBindJSON(&user); err == nil {
		// user := models.User{
		// 	Name:           json.Name,
		// 	Surname:        json.Surname,
		// 	Patronymic:     json.Patronymic,
		// 	PassportNumber: json.PassportNumber,
		// 	Address:        json.Address,

		// }

		err := h.userService.CreateUser(user)
		fmt.Println("handlers", user.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "created", "user": user})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	// Получение параметров запроса
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	// Сбор фильтров из параметров запроса
	filters := make(map[string]string)
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}
	if passportNumber := c.Query("passportNumber"); passportNumber != "" {
		filters["passportNumber"] = passportNumber
	}

	// Получение пользователей из сервиса
	users, total, err := h.userService.GetUsers(page, limit, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Формирование ответа с пагинацией
	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"total": total,
		"data":  users,
	})
}

func (h *UserHandler) GetWorklogs(c *gin.Context) {
	// Реализация логики получения трудозатрат пользователя
	userId := c.Param("userId")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	c.JSON(http.StatusOK, gin.H{
		"userId":    userId,
		"startDate": startDate,
		"endDate":   endDate,
		"worklogs": []gin.H{
			{"taskId": 101, "taskName": "Task 1", "hours": 5, "minutes": 30},
			{"taskId": 102, "taskName": "Task 2", "hours": 3, "minutes": 45},
		},
	})
}

func (h *UserHandler) StartTask(c *gin.Context) {
	// Реализация логики начала отсчета времени по задаче
	c.JSON(http.StatusOK, gin.H{"status": "started", "timestamp": "2023-01-01T12:00:00Z"})
}

func (h *UserHandler) StopTask(c *gin.Context) {
	// Реализация логики окончания отсчета времени по задаче
	c.JSON(http.StatusOK, gin.H{"status": "stopped", "timestamp": "2023-01-01T14:30:00Z"})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userId := c.Param("userId")

	// Логирование начала процесса удаления пользователя
	logger.PrintDebug("Attempting to delete user", map[string]any{"userId": userId})

	// Вызов метода удаления пользователя из сервиса
	if err := h.userService.DeleteUser(userId); err != nil {
		logger.PrintError("Failed to delete user", map[string]any{"userId": userId, "error": err.Error()})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Логирование успешного удаления пользователя
	logger.PrintInfo("User deleted successfully", map[string]any{"userId": userId})

	// Возврат успешного ответа
	c.JSON(http.StatusOK, gin.H{"status": "deleted", "userId": userId})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	// Реализация логики обновления данных пользователя
	userId := c.Param("userId")
	var json struct {
		Name           string `json:"name"`
		PassportNumber string `json:"passportNumber"`
	}
	if err := c.ShouldBindJSON(&json); err == nil {
		c.JSON(http.StatusOK, gin.H{"status": "updated", "user": gin.H{"id": userId, "name": json.Name, "passportNumber": json.PassportNumber}})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
