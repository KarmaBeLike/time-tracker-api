package handlers

import (
	"errors"
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
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Routes(router *gin.Engine, cfg *config.Config) {
	user := router.Group("/users")
	{
		user.POST("/", h.CreateUser)
		user.GET("/", h.GetUsers)
		user.DELETE("/:userId", h.DeleteUser)
		user.PUT("/:userId", h.UpdateUser)

	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	logger.PrintInfo("Handling CreateUser request", nil)

	var json struct {
		PassportNumber string `json:"passportNumber"`
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := h.userService.CreateUser(json.PassportNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "created", "user": user})
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	// Получение параметров запроса
	logger.PrintInfo("Handling GetUsers request", nil)

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

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userId := c.Param("userId")

	// Логирование начала процесса удаления пользователя
	logger.PrintDebug("Attempting to delete user", map[string]any{"userId": userId})

	// Вызов метода удаления пользователя из сервиса
	if err := h.userService.DeleteUser(userId); err != nil {
		logger.PrintError(errors.New("Failed to delete user"), map[string]any{"userId": userId, "error": err.Error()})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Логирование успешного удаления пользователя
	logger.PrintInfo("User deleted successfully", map[string]any{"userId": userId})

	// Возврат успешного ответа
	c.JSON(http.StatusOK, gin.H{"status": "deleted", "userId": userId})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var json struct {
		Name           string `json:"name"`
		Surname        string `json:"surname"`
		Patronymic     string `json:"patronymic"`
		PassportNumber string `json:"passportNumber"`
		Address        string `json:"address"`
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &models.User{
		ID:             userId,
		Name:           json.Name,
		Surname:        json.Surname,
		Patronymic:     json.Patronymic,
		PassportNumber: json.PassportNumber,
		Address:        json.Address,
	}

	if err := h.userService.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated", "user": user})
}
