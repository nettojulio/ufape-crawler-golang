package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nettojulio/ufape-crawler-golang/server/handlers"

	"github.com/go-playground/validator/v10"
	"github.com/nettojulio/ufape-crawler-golang/config"
)

func New(cfg *config.Config) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())

	e.Validator = &CustomValidator{validator: validator.New()}

	e.POST("/", handlers.CrawlerHandler)
	e.GET("/", handlers.HealthCheckHandler(cfg.Version))

	return e
}
