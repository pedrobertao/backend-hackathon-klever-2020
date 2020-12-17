package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/joho/godotenv"
	"github.com/pedrobertao/backend-hackathon-klever-2020/database"
	"github.com/pedrobertao/backend-hackathon-klever-2020/models"
	"github.com/pedrobertao/backend-hackathon-klever-2020/sms"
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
		log.Fatal("Port not set")
	}
	router.Run(":" + port)
}

func main() {
	zapConfig := zap.NewDevelopmentConfig()
	logger, err := zapConfig.Build()
	if err != nil {
		log.Fatal("Error to init zap global logger")
	}
	zap.ReplaceGlobals(logger)

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	zap.L().Info("Env loaded")

	if err := database.Connect(); err != nil {
		log.Fatal("Error connecting to database")
	}
	defer database.Stop()
	zap.L().Info("Database connected")

	serve()

	sms.SendSMS(models.SMS{
		To:   "+5585999263009",
		From: "+12517149048",
		Body: "You have received 1 BTC.",
	})
}
