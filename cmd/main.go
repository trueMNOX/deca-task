package main

import (
	"deca-task/internal/auth"
	"deca-task/internal/config"
	"deca-task/internal/database"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	db := database.InitDB(cfg)
	redisClient := database.InitRedis(cfg)

	authRepo := auth.NewAuthRepository(db, redisClient)
	authService := auth.NewAuthService(authRepo)
	authHandler := auth.NewAuthHandler(authService)

	r := gin.Default()

	v1 := r.Group("/api/v1")
	authHandler.AuthRoute(v1)

	r.Run(":8080")
}
