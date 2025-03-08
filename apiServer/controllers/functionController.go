package controllers

import (
	"faas/models"
	"faas/services/functions"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FunctionController struct {
	service *functions.FunctionService
}

func NewFunctionController(service *functions.FunctionService) *FunctionController {
	return &FunctionController{
		service: service,
	}
}

func (f *FunctionController) RegisterFunction(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	var req functions.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	funcID, err := f.service.RegisterFunction(c.Request.Context(), user.(models.User), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, functions.RegisterResponse{FunctionIdentifier: funcID})
}

func (f *FunctionController) DeleteFunction(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	functionID := c.Param("id")
	err := f.service.DeleteFunction(c.Request.Context(), user.(models.User), functionID)

	switch err {
	case nil:
		// No hay error
		c.JSON(http.StatusOK, gin.H{"message": "Function deleted successfully"})
	case functions.ErrFunctionNotFound:
		// La función no se ha encontrado en el kvStore
		c.JSON(http.StatusNotFound, gin.H{"error": "Function not found"})
	case functions.ErrUnauthorized:
		// El usuario no tiene permisos para borrar esa función
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to delete this function"})
	default:
		// Otro tipo de error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete function"})
	}
}

func (f *FunctionController) GetFunctions(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	functions, err := f.service.GetUserFunctions(c.Request.Context(), user.(models.User).Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve functions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"functions": functions})
}
