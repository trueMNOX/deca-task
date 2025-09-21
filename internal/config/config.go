package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	JWTSecret   string
	JWTExpireIn int
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("[WARN] .env file not found, using system environment variables")
	}

	redisDB, err := strconv.Atoi(getRequiredEnv("REDIS_DB"))
	if err != nil {
		log.Fatalf("Invalid REDIS_DB value: %v", err)
	}

	jwtExpireIn, err := strconv.Atoi(getRequiredEnv("JWT_EXPIRE_IN"))
	if err != nil {
		log.Fatalf("Invalid JWT_EXPIRE_IN value: %v", err)
	}

	return &Config{
		AppPort: getRequiredEnv("APP_PORT"),

		PostgresHost:     getRequiredEnv("POSTGRES_HOST"),
		PostgresPort:     getRequiredEnv("POSTGRES_PORT"),
		PostgresUser:     getRequiredEnv("POSTGRES_USER"),
		PostgresPassword: getRequiredEnv("POSTGRES_PASSWORD"),
		PostgresDB:       getRequiredEnv("POSTGRES_DB"),

		RedisHost:     getRequiredEnv("REDIS_HOST"),
		RedisPort:     getRequiredEnv("REDIS_PORT"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,

		JWTSecret:   getRequiredEnv("JWT_SECRET"),
		JWTExpireIn: jwtExpireIn,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getRequiredEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Fatalf("Required environment variable %s is not set", key)
	return ""
}
