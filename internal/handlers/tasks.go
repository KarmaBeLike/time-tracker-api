package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/KarmaBeLike/time-tracker-api/config"
	"github.com/KarmaBeLike/time-tracker-api/internal/service"
	"github.com/KarmaBeLike/time-tracker-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (t *TaskHandler) Routes(router *gin.Engine, cfg *config.Config) {
	task := router.Group("/tasks")
	{
		task.GET("/:userId/worklogs", t.GetWorklogs)
		task.POST("/:userId/tasks/:taskId/start", t.StartTask)
		task.POST("/:userId/tasks/:taskId/stop", t.StopTask)
	}
}

func (h *TaskHandler) GetWorklogs(c *gin.Context) {
	userIdStr := c.Param("userId")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		logger.PrintError(errors.New("invalid user ID"), map[string]any{"userId": userIdStr})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	worklogs, err := h.taskService.GetWorklogs(userId, startDate, endDate)
	if err != nil {
		logger.PrintError(errors.New("failed to get worklogs"), map[string]any{"userId": userId, "startDate": startDate, "endDate": endDate, "error": err.Error()})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get worklogs"})
		return
	}

	logger.PrintInfo("Worklogs retrieved successfully", map[string]any{"userId": userId, "startDate": startDate, "endDate": endDate})

	c.JSON(http.StatusOK, gin.H{
		"userId":    userId,
		"startDate": startDate,
		"endDate":   endDate,
		"worklogs":  worklogs,
	})
}

func (t *TaskHandler) StartTask(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	taskId, err := strconv.Atoi(c.Param("taskId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	err = t.taskService.StartTask(userId, taskId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "started", "timestamp": time.Now().UTC().Format(time.RFC3339)})
}

func (t *TaskHandler) StopTask(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	taskId, err := strconv.Atoi(c.Param("taskId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	err = t.taskService.StopTask(userId, taskId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "stopped", "timestamp": time.Now().UTC().Format(time.RFC3339)})
}
