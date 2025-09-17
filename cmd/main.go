package main

import (
	"deca-task/internal/auth"
	"deca-task/internal/config"
	"deca-task/internal/database"
	"deca-task/internal/middleware"
	"deca-task/internal/user"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	db := database.InitDB(cfg)
	redisClient := database.InitRedis(cfg)

	authRepo := auth.NewAuthRepository(db, redisClient)
	authService := auth.NewAuthService(authRepo)
	authHandler := auth.NewAuthHandler(authService)
	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	r := gin.Default()

	v1 := r.Group("/api/v1")
	authHandler.AuthRoute(v1)
	v2 := r.Group("/api/v2", middleware.AuthModdleware())
	userHandler.UsersRoute(v2)

	r.Run(":8080")
}
