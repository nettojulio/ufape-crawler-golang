package http

import (
	"github.com/labstack/echo/v4"
	"github.com/nettojulio/ufape-crawler-golang/internal/crawler"
)

// HealthCheckHandler godoc
// @Summary      Verifica a saúde da API
// @Description  Retorna o status "OK" e a versão atual da aplicação.
// @Tags         Health
// @Produce      json
// @Success      200  {object} crawler.APIHealth "API está saudável"
// @Router       / [get]
func HealthCheckHandler(version string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(200, crawler.APIHealth{
			Status:  "OK",
			Version: version,
		})
	}
}
