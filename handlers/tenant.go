package handlers

import (
	"context"
	"encoding/json"
	"jatis_mobile_api/database"
	"jatis_mobile_api/logs"
	"jatis_mobile_api/models"
	"jatis_mobile_api/rabbitmq"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func CreateTenantHandler(c echo.Context) error {
	logger := c.Get("logger").(*logrus.Logger)
	var tenant models.Tenant

	if err := c.Bind(&tenant); err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to bind request for tenant", struct{ Error error }{Error: err})
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	db := database.GetDB()

	var existingTenant models.Tenant
	err := db.QueryRow(context.Background(), "SELECT id FROM tenants WHERE name = $1", tenant.Name).Scan(&existingTenant.ID)
	if err == nil {
		logs.LogWithFields(logger, logrus.WarnLevel, "Attempted to create tenant but it already exists", struct{ TenantName string }{TenantName: tenant.Name})
		return c.JSON(http.StatusConflict, "Tenant already exists")
	}

	if err := models.CreateTenant(db, &tenant); err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to create tenant", struct {
			TenantName string
			Error      error
		}{TenantName: tenant.Name, Error: err})
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := rabbitmq.DeclareQueue(tenant.Name); err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to queue tenant created to RabbitMQ", struct {
			TenantName string
			Error      error
		}{TenantName: tenant.Name, Error: err})
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := rabbitmq.BindQueue(tenant.Name, "amq.direct", tenant.Name); err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to bind queue tenant to RabbitMQ", struct {
			TenantName string
			Error      error
		}{TenantName: tenant.Name, Error: err})
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	successMessage := map[string]string{"message": "Tenant created successfully"}
	successBody, err := json.Marshal(successMessage)
	if err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to marshal success message for RabbitMQ", struct {
			TenantName string
			Error      error
		}{TenantName: tenant.Name, Error: err})
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := rabbitmq.PublishMessage("amq.direct", tenant.Name, successBody); err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to publish tenant created message to RabbitMQ", struct {
			TenantName string
			Error      error
		}{TenantName: tenant.Name, Error: err})
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	logger.WithField("tenant_name", tenant.Name).Info("Tenant created successfully")
	return c.JSON(http.StatusCreated, tenant)
}

func DeleteTenantHandler(c echo.Context) error {
	logger := c.Get("logger").(*logrus.Logger)
	tenantIDStr := c.Param("id")

	tenantID, err := strconv.Atoi(tenantIDStr)
	if err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Invalid tenant ID", struct{ TenantID string }{TenantID: tenantIDStr})
		return c.JSON(http.StatusBadRequest, "Invalid tenant ID")
	}

	db := database.GetDB()

	var tenant models.Tenant
	err = db.QueryRow(context.Background(), "SELECT id, name FROM tenants WHERE id = $1", tenantID).Scan(&tenant.ID, &tenant.Name)
	if err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to retrieve tenant", struct{ TenantID int }{TenantID: tenantID})
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := models.SoftDeleteTenant(db, tenantID); err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to soft delete tenant", struct{ TenantName string }{TenantName: tenant.Name})
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	channel, err := rabbitmq.GetChannel()
	if err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to get RabbitMQ channel", struct{ TenantName string }{TenantName: tenant.Name})
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer channel.Close()

	queueName := tenant.Name

	messageCount, err := channel.QueueDelete(queueName, false, false, false)
	if err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to delete queue", struct{ QueueName string }{QueueName: queueName})
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	logs.LogWithFields(logger, logrus.InfoLevel, "Queue deleted successfully", struct {
		QueueName    string
		MessageCount int
	}{QueueName: queueName, MessageCount: messageCount})

	logs.LogWithFields(logger, logrus.InfoLevel, "Tenant deleted successfully", struct{ TenantName string }{TenantName: tenant.Name})
	return c.JSON(http.StatusOK, "Tenant deleted successfully")
}
