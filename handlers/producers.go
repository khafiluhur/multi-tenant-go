package handlers

import (
	"encoding/json"
	"jatis_mobile_api/logs"
	"jatis_mobile_api/rabbitmq"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func ProducerHandler(c echo.Context) error {
	logger := c.Get("logger").(*logrus.Logger)

	tenantName := c.Request().Header.Get("x-tenant-name")
	if tenantName == "" {
		logs.LogWithFields(logger, logrus.ErrorLevel, "x-tenant-name header is required", struct{}{})
		return c.JSON(http.StatusBadRequest, "x-tenant-name header is required")
	}

	var requestBody map[string]interface{}
	if err := c.Bind(&requestBody); err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to bind request body", struct{ Error string }{Error: err.Error()})
		return c.JSON(http.StatusBadRequest, "Invalid request body")
	}

	messageJSON, err := json.Marshal(requestBody)
	if err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to marshal message to JSON", struct{ Error string }{Error: err.Error()})
		return c.JSON(http.StatusInternalServerError, "Failed to process message")
	}

	queueName := tenantName

	if err := rabbitmq.PublishMessage("amq.direct", queueName, messageJSON); err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to start RabbitMQ consumer", struct{ QueueName string }{QueueName: queueName})
		return c.JSON(http.StatusInternalServerError, "Failed to start consumer")
	}

	response := map[string]interface{}{
		"tenant_name": tenantName,
		"message":     requestBody,
	}

	return c.JSON(http.StatusOK, response)
}
