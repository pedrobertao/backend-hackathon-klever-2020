package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/joho/godotenv"
	"github.com/pedrobertao/backend-hackathon-klever-2020/database"
	"github.com/pedrobertao/backend-hackathon-klever-2020/models"
	"github.com/pedrobertao/backend-hackathon-klever-2020/sms"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	MainAddress string `bson:"mainAddress" json:"mainAddress" binding:"required"`
	Username    string `bson:"username" json:"username" binding:"required"`
	Email       string `bson:"email" json:"email" binding:"required"`
	Phone       string `bson:"phone" json:"phone" binding:"required"`
}

type SMS struct {
	From   string  `json:"from" binding:"required"`
	To     string  `json:"to" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
	Coin   string  `json:"coin" binding:"required"`
}

func serve() {
	router := gin.New()

	router.GET("/user/:username", func(c *gin.Context) {
		var getRequest struct {
			Username string `json:"username" uri:"username" binding:"required"`
		}
		if err := c.ShouldBindUri(&getRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		filter := bson.M{"username": getRequest.Username}
		options := options.FindOne()

		var u User
		if err := database.UsersCollection.FindOne(c, filter, options).Decode(&u); err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "User not registered",
				})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, u)
	})

	router.POST("/user/phone", func(c *gin.Context) {
		var phoneVerification struct {
			Phone string `json:"phone" uri:"phone" binding:"required"`
			Code  string `json:"code" uri:"code" binding:""`
		}

		if err := c.ShouldBindJSON(&phoneVerification); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		if phoneVerification.Code != "" {
			err := sms.VerifyCodeSMS(phoneVerification.Phone, phoneVerification.Code)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "Phone verified",
			})
			return
		}

		err := sms.SendVerifySMS(phoneVerification.Phone)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Verification code sent",
		})
		return
	})

	router.PUT("/user", func(c *gin.Context) {
		var userRequest User
		if err := c.ShouldBindJSON(&userRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		_, err := database.UsersCollection.InsertOne(c, userRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User registered"})
	})

	router.POST("/sms/transaction", func(c *gin.Context) {
		var smsRequest SMS
		if err := c.ShouldBindJSON(&smsRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		phone := ""
		// TODO - Retrieve from DATABASE
		if smsRequest.To == "bertao" {
			phone = "+5531996139388"
		} else if smsRequest.To == "roney" {
			phone = "+5585999263009"
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Username not found",
			})
			return
		}

		// TODO - Need to check from
		if err := sms.SendSMS(models.SMS{
			To:   phone,
			From: "+12517149048",
			Body: fmt.Sprintf("You have received %f %s from %s",
				smsRequest.Amount, smsRequest.Coin, smsRequest.From),
		}); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Text sent"})
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "Klever ID is live",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		zap.L().Fatal("Port not set")
	}
	gin.SetMode(gin.ReleaseMode)

	if err := router.Run(":" + port); err != nil {
		zap.L().Fatal("Router fail", zap.Error(err))
	}

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

	if err := sms.Config(); err != nil {
		zap.L().Fatal("Error loading SMS config")
	}

	serve()
}
