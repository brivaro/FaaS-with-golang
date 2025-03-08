package controllers

import (
	"faas/models"
	"faas/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	repo *repository.UserRepository
}

func NewUserController(repo *repository.UserRepository) *UserController {
	return &UserController{
		repo: repo,
	}
}

func (u *UserController) GetAllUsers(c *gin.Context) {
	user, _ := c.Get("user")
	if user.(models.User).Role != "admin" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	users, err := u.repo.GetAllUsers()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, users)
}
