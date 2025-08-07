package handlers

import "github.com/labstack/echo/v4"

func HealthCheckHandler(version string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status":  "OK",
			"version": version,
		})
	}
}
