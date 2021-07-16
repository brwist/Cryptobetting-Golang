package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RootRoute (c echo.Context) error {
	return c.HTML(http.StatusOK, "Welcome Gambling!")
}
