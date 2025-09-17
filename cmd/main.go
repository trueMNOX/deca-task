// @title OTP Auth Service API
// @version 1.0
// @description Backend service in Golang for OTP-based login & user management
// @contact.name Mehdi Dev
// @host localhost:8080
// @BasePath /
package main

import (
	_ "deca-task/docs"
	"deca-task/internal/auth"
	"deca-task/internal/config"
	"deca-task/internal/database"
	"deca-task/internal/middleware"
	"deca-task/internal/user"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
}
