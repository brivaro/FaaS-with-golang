package controllers

import (
	"faas/models"
	"faas/services/executor"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExecutorController struct {
	service *executor.ExecutorService
}

func NewExecutorController(executorService *executor.ExecutorService) *ExecutorController {
	return &ExecutorController{
		service: executorService,
	}
}

func (e *ExecutorController) ExecuteFunction(c *gin.Context) {
	user, exists := c.Get("user")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	var req executor.ExecuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	task, err := e.service.ExecuteFunction(c.Request.Context(), user.(models.User), req)
	if err != nil {
		switch err {
		case executor.ErrUnauthorized:
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to execute function: %s", err.Error())})
		}
		return
	}

	result, err := e.service.GetResult(c.Request.Context(), task)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.Error != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, result)
}
