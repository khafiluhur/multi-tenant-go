package main

import (
	"fmt"
	"jatis_mobile_api/config"
	"jatis_mobile_api/database"
	"jatis_mobile_api/logs"
	"jatis_mobile_api/middleware"
	"jatis_mobile_api/migrations"
	"jatis_mobile_api/rabbitmq"
	"jatis_mobile_api/routes"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var logger = logs.SetupLogger()

func main() {
	logger.Info("Loading configuration...")
	cfg, err := config.LoadConfig()
	if err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Could not load config", struct{ Error error }{Error: err})
		return
	}

	logger.Info("Connecting to database...")
	if err := database.ConnectDB(cfg.PostgresURL); err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Could not connect to database", struct {
			PostgresURL string
			Error       error
		}{PostgresURL: cfg.PostgresURL, Error: err})
		return
	}

	logger.Info("Running migrations...")
	db := database.GetDB()

	if err := migrations.CreateTenantsTable(db, logger); err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to create tenants table", struct{ Error error }{Error: err})
		return
	}

	logger.Info("Connecting to RabbitMQ...")
	if err := rabbitmq.ConnectRabbitMQ(cfg.RabbitMQURL); err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Could not connect to RabbitMQ", struct {
			RabbitMQURL string
			Error       error
		}{RabbitMQURL: cfg.RabbitMQURL, Error: err})
		return
	}

	go monitorRabbitMQConnection(cfg.RabbitMQURL)

	e := echo.New()
	e.Use(middleware.LoggerMiddleware)
	e.Use(middleware.PerformanceLogger(logger))
	routes.RegisterTenantRoutes(e)

	address := fmt.Sprintf(":%d", cfg.PORT)
	logs.LogWithFields(logger, logrus.InfoLevel, "Starting server", struct{ Port int }{Port: cfg.PORT})
	e.Logger.Fatal(e.Start(address))
}

func monitorRabbitMQConnection(url string) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	isActive := false

	for range ticker.C {
		if rabbitmq.IsClosed() {
			if isActive {
				logs.LogWithFields(logger, logrus.ErrorLevel, "RabbitMQ connection is closed!", struct{}{})
				isActive = false
			}
			if err := rabbitmq.ConnectRabbitMQ(url); err != nil {
				logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to reconnect to RabbitMQ", struct{ RabbitMQURL string }{RabbitMQURL: url})
			} else {
				logs.LogWithFields(logger, logrus.InfoLevel, "Reconnected to RabbitMQ successfully", struct{}{})
				isActive = true
			}
		} else {
			if !isActive {
				logs.LogWithFields(logger, logrus.InfoLevel, "RabbitMQ connection is active", struct{}{})
				isActive = true
			}
		}
	}
}
