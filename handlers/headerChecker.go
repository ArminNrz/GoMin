package handlers

import (
	"GoMin/config"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CheckHeader(c echo.Context) error {
	apiKey := c.Request().Header.Get("X-API-KEY")
	serviceName := c.Request().Header.Get("Service-name")

	if apiKey == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "X-API-KEY header is missing")
	}

	if serviceName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Service-name header is missing")
	}

	if key, exist := config.AppConfig.API.Keys[serviceName]; !exist {
		return echo.NewHTTPError(http.StatusNotFound, "Service not found")
	} else if apiKey != key {
		return echo.NewHTTPError(http.StatusUnauthorized, "Service key is invalid")
	}

	return nil
}
