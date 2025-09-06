package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nettojulio/ufape-crawler-golang/config"
	"github.com/nettojulio/ufape-crawler-golang/server/handlers"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/nettojulio/ufape-crawler-golang/docs"
)

func New(cfg *config.Config) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Validator = &CustomValidator{validator: validator.New()}

	e.POST("/", handlers.CrawlerHandler)
	e.GET("/", handlers.HealthCheckHandler(cfg.Version))

	return e
}
