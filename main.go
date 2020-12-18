package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

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

		var u models.User
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

	router.GET("/user", func(c *gin.Context) {
		var getRequest struct {
			Search string `json:"search" form:"search" binding:""`
		}
		if err := c.ShouldBindQuery(&getRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		filter := bson.M{
			"$or": []bson.M{
				{"username": getRequest.Search},
				{"$and": []bson.M{
					{"phone.phone": getRequest.Search},
					{"phone.status": models.Active},
				}},
				{"$and": []bson.M{
					{"email.email": getRequest.Search},
					{"email.status": models.Active},
				}},
			},
		}
		options := options.FindOne()
		var u models.User
		if err := database.UsersCollection.FindOne(c, filter, options).Decode(&u); err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"addresses": u.Addresses,
		})

	})

	router.POST("/user/phone", func(c *gin.Context) {
		var phoneVerification struct {
			Username string `json:"username" uri:"username" binding:"required"`
			Code     string `json:"code" uri:"code" binding:""`
		}

		if err := c.ShouldBindJSON(&phoneVerification); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		var user models.User
		if err := database.UsersCollection.FindOne(c, bson.M{"username": phoneVerification.Username}).Decode(&user); err != nil {
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

		if user.Phone.Status == models.Active {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Phone already active",
			})
			return
		}

		phoneToVerify := "+" + user.Phone.Phone
		if phoneVerification.Code != "" {
			err := sms.VerifyCodeSMS(phoneToVerify, phoneVerification.Code)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": err.Error(),
				})
				return
			}

			filter := bson.M{"username": phoneVerification.Username}
			updated := bson.M{"$set": bson.M{"phone.status": models.Active}}

			var user models.User
			if err := database.UsersCollection.FindOneAndUpdate(c, filter, updated).Decode(&user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Phone verified"})
			return
		}

		if err := sms.SendVerifySMS(phoneToVerify); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Verification code sent"})
		return
	})

	router.PUT("/user", func(c *gin.Context) {
		var userRequest struct {
			Addresses   []string `json:"addresses" binding:"required,min=1"`
			MainAddress string   `json:"mainAddress" binding:"required"`
			Username    string   `json:"username" binding:"required"`
			Email       string   `json:"email"`
			Phone       string   `json:"phone"`
		}

		if err := c.ShouldBindJSON(&userRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		filter := bson.M{
			"$or": []bson.M{
				{"username": userRequest.Username},
				{"$and": []bson.M{
					{"phone.phone": bson.M{"$exists": true}},
					{"phone.phone": userRequest.Phone},
				}},
				{"$and": []bson.M{
					{"email.email": bson.M{"$exists": true}},
					{"email.email": userRequest.Phone},
				}},
			},
		}

		var user models.User
		if err := database.UsersCollection.FindOne(c, filter).Decode(&user); err != nil {
			if err == mongo.ErrNoDocuments {
				_, err = database.UsersCollection.InsertOne(c, models.User{
					Addresses:   userRequest.Addresses,
					MainAddress: userRequest.MainAddress,
					Username:    userRequest.Username,
					Email: models.UserEmail{
						Email:  userRequest.Email,
						Status: models.Inactive,
					},
					Phone: models.UserPhone{
						Phone:  userRequest.Phone,
						Status: models.Inactive,
					},
					UpdatedAt: time.Now(),
				})
				if err == nil {
					c.JSON(http.StatusOK, gin.H{"message": "User registered"})
					return
				}
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User Already Exists",
		})
		return
	})

	router.POST("/sms/transaction", func(c *gin.Context) {
		var smsRequest models.SMSRequest
		if err := c.ShouldBindJSON(&smsRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		phone := ""
		// TODO - Retrieve from DATABASE
		if smsRequest.To == "bertao" {
			phone = "+" + "5531996139388"
		} else if smsRequest.To == "roney" {
			phone = "+" + "5585999263009"
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Username not found",
			})
			return
		}

		// TODO - Need to check from
		if err := sms.SendSMS(models.SMS{
			To:   phone,
			From: "+18058645005",
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

	//sms.SendSMS(models.SMS{
	//To:   "+5585999263009",
	//From: "+18058645005",
	//Body: "You have received 1 BTC.",
	//})
}
