package controllers

import (
	"faas/dataSource"
	"faas/models"
	"faas/services/auth"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	authService *auth.AuthService
}

func NewMiddleware(authService *auth.AuthService) *Middleware {
	return &Middleware{
		authService: authService,
	}
}

func (m *Middleware) RequireAuth(c *gin.Context) {
	// Get the cookie off req
	var tokenString string
	var err error

	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	} else {
		// Intentar obtener el token de la cookie Authorization
		tokenString, err = c.Cookie("Authorization")
		if err != nil {
			// Si no está presente en ningún lugar, responder con Unauthorized
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}

	claims, err := m.authService.Validate(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	if claims.ExpiresAt.Compare(time.Now()) <= 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var user models.User
	user, err = dataSource.GetUserByUsername(claims.Username)

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("user", user)
	c.Next()
}
