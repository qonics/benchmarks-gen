package config

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func InitializeConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file, system will load default value")
	}
	if os.Getenv("APP_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}
