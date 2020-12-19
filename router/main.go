package router

import (
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Serve() {
	router := gin.New()

	router.GET("/user/:address", GetUserByAddress)

	router.GET("/user", GetUser)

	router.POST("/user/phone", PhoneVerify)

	router.PUT("/user", CreateUser)

	router.POST("/sms/transaction", SmsTransaction)

	router.GET("/", Home)

	port := os.Getenv("PORT")
	if port == "" {
		zap.L().Fatal("Port not set")
	}
	gin.SetMode(gin.ReleaseMode)

	if err := router.Run(":" + port); err != nil {
		zap.L().Fatal("Router fail", zap.Error(err))
	}

}
