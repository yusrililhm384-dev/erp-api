package routes

import (
	"log"
	"os"
	"time"

	midd "enterprise_resource_planning/internal/platform/middleware"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

// @title ERP OpenAPI
// @version 1.0
// @description This is an api documentation for XiaoERP.
// @termsOfService http://terms.xiaoching.com

// @contact.name Xiaoching Support
// @contact.url http://developer.xiaoching.com
// @contact.email support@xiaoching.com

// @host localhost
// @BasePath /v1
func New() {
	e := echo.New()

	// login purpose
	e.Use(midd.RateLimit(3, 5, 2*time.Minute))

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	if err := e.Start(":" + os.Getenv("APP_PORT")); err != nil {
		log.Fatal(err.Error())
	}
}
