package http

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/nettojulio/ufape-crawler-golang/docs"
	"github.com/nettojulio/ufape-crawler-golang/internal/config"
)

func registerRoutes(e *echo.Echo, cfg *config.Config, crawlerHandler *CrawlerHandler) {
	healthCheckRoutes(e, cfg)
	crawlerRoutes(e, crawlerHandler)
	swaggerRoutes(e)
}

func healthCheckRoutes(e *echo.Echo, cfg *config.Config) {
	e.GET("/", HealthCheckHandler(cfg.Version))
}

func crawlerRoutes(e *echo.Echo, h *CrawlerHandler) {
	e.POST("/", h.HandleCrawl)
}

func swaggerRoutes(e *echo.Echo) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
