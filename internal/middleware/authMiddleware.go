package middleware

import (
	"deca-task/internal/auth/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthModdleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			return
		}
		token, err := jwt.ParseToken(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		userId, err := jwt.GetUserId(*token)
		c.Set("user_id", userId)
		c.Next()
	}
}
