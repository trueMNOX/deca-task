package database

import (
	"deca-task/internal/config"
	"deca-task/internal/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) *gorm.DB {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDB,
		cfg.PostgresPort,
	)

	var err error

	db , err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to Postgres: %v", err)
	}
	log.Println("✅ Postgres connected")

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("migrate is failed")
	}
	log.Println("✅ migrated is ok")

	DB = db

	return db
}
