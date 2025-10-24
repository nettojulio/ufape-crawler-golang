package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/nettojulio/ufape-crawler-golang/internal/config"
	"github.com/nettojulio/ufape-crawler-golang/internal/crawler"
)

// NewServer agora aceita o crawler.Service como uma dependência.
func NewServer(cfg *config.Config, crawlerService *crawler.Service) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Validator = &CustomValidator{validator: validator.New()}

	crawlerHandler := NewCrawlerHandler(crawlerService)

	registerRoutes(e, cfg, crawlerHandler)

	return e
}
