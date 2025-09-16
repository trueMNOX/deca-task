package main

import (
	"deca-task/internal/config"
	"deca-task/internal/database"
)

func main(){
	cfg := config.LoadConfig()
	database.InitDB(cfg)
	database.InitRedis(cfg)
}
