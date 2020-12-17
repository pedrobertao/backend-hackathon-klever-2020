package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "3030"
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/", func(c *gin.Context) {
		// c.HTML(http.StatusOK, "Klever ID is live", nil)
		c.JSON(http.StatusOK, gin.H{
			"status": "Klever ID is live",
		})
	})

	router.Run(":" + port)
}
