package handlers

import (
	"net/http"
	"strconv"

	"github.com/KarmaBeLike/time-tracker-api/config"
	"github.com/KarmaBeLike/time-tracker-api/internal/service"
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
		user.GET("/", h.GetUsers)
		user.GET("/:userId/worklogs", h.GetWorklogs)
		user.POST("/:userId/tasks/:taskId/start", h.StartTask)
		user.POST("/:userId/tasks/:taskId/stop", h.StopTask)
		user.DELETE("/:userId", h.DeleteUser)
		user.PUT("/:userId", h.UpdateUser)
		user.POST("/", h.CreateUser)
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
	// Реализация логики удаления пользователя
	userId := c.Param("userId")
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

func (h *UserHandler) CreateUser(c *gin.Context) {
	// Реализация логики создания нового пользователя
	var json struct {
		PassportNumber string `json:"passportNumber"`
	}
	if err := c.ShouldBindJSON(&json); err == nil {
		c.JSON(http.StatusOK, gin.H{"status": "created", "user": gin.H{"id": 2, "passportNumber": json.PassportNumber}})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
