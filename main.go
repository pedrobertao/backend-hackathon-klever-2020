package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/joho/godotenv"
	"github.com/pedrobertao/backend-hackathon-klever-2020/database"
	"go.uber.org/zap"
)

func serve() {
	router := gin.New()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "Klever ID is live",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		zap.L().Fatal("Port not set")
	}
	router.Run(":" + port)
}

func main() {
	zapConfig := zap.NewDevelopmentConfig()
	logger, err := zapConfig.Build()
	if err != nil {
		zap.L().Fatal("Error to init zap global logger")
	}
	zap.ReplaceGlobals(logger)

	err = godotenv.Load()
	if err != nil {
		zap.L().Warn(".env not found, using os variables")
	} else {
		zap.L().Info("Env loaded")
	}

	if err := database.Connect(); err != nil {
		zap.L().Fatal("Error connecting to database")
	}
	defer database.Stop()
	zap.L().Info("Database connected")

	serve()
}
