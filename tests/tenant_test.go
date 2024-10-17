package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"jatis_mobile_api/database"
	"jatis_mobile_api/handlers"
	"jatis_mobile_api/logs"
	"jatis_mobile_api/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupEcho() *echo.Echo {
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger := logs.SetupLogger()
			c.Set("logger", logger)
			return next(c)
		}
	})
	return e
}

func setupDatabase() {
	if err := database.ConnectDB("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"); err != nil {
		panic(err)
	}
}

func TestCreateTenantHandler(t *testing.T) {
	setupDatabase()
	defer database.Close()

	e := setupEcho()

	tenant := models.Tenant{Name: "Test Tenant"}
	tenantJSON, _ := json.Marshal(tenant)

	req := httptest.NewRequest(http.MethodPost, "/tenants", bytes.NewBuffer(tenantJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handlers.CreateTenantHandler(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		var createdTenant models.Tenant
		err := json.Unmarshal(rec.Body.Bytes(), &createdTenant)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		assert.NotEqual(t, 0, createdTenant.ID)
	}
}

func TestDeleteTenantHandler(t *testing.T) {
	setupDatabase()
	defer database.Close()

	e := setupEcho()

	tenant := models.Tenant{Name: "Test Tenant To Delete"}
	tenantJSON, _ := json.Marshal(tenant)

	createReq := httptest.NewRequest(http.MethodPost, "/tenants", bytes.NewBuffer(tenantJSON))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	createCtx := e.NewContext(createReq, createRec)

	handlers.CreateTenantHandler(createCtx)

	var createdTenant models.Tenant
	json.Unmarshal(createRec.Body.Bytes(), &createdTenant)

	req := httptest.NewRequest(http.MethodDelete, "/tenants/"+strconv.Itoa(createdTenant.ID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handlers.DeleteTenantHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Tenant deleted successfully", rec.Body.String())
	}
}
