package main

import (
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/joho/godotenv"
	"github.com/pedrobertao/backend-hackathon-klever-2020/database"
	"github.com/pedrobertao/backend-hackathon-klever-2020/router"
	"github.com/pedrobertao/backend-hackathon-klever-2020/sms"
	"go.uber.org/zap"
)

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

	router.Serve()

}
