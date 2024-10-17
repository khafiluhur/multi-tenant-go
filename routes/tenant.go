package routes

import (
	"jatis_mobile_api/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterTenantRoutes(e *echo.Echo) {
	e.POST("/tenants", handlers.CreateTenantHandler)
	e.DELETE("/tenants/:id", handlers.DeleteTenantHandler)
	e.GET("/consumers", handlers.ConsumerHandler)
	e.POST("/producers", handlers.ProducerHandler)
}
