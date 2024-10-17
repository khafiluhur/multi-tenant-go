package handlers

import (
	"jatis_mobile_api/logs"
	"jatis_mobile_api/rabbitmq"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func ConsumerHandler(c echo.Context) error {
	logger := c.Get("logger").(*logrus.Logger)

	tenantName := c.Request().Header.Get("x-tenant-name")
	if tenantName == "" {
		logs.LogWithFields(logger, logrus.ErrorLevel, "x-tenant-name header is required", struct{}{})
		return c.JSON(http.StatusBadRequest, "x-tenant-name header is required")
	}

	queueName := tenantName

	if err := rabbitmq.ConsumeMessages(queueName); err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to start RabbitMQ consumer", struct{ QueueName string }{QueueName: queueName})
		return c.JSON(http.StatusInternalServerError, "Failed to start consumer")
	}

	return c.JSON(http.StatusOK, "RabbitMQ consumer started successfully and message published")
}
