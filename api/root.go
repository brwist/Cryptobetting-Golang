package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var RootGreeting = "Welcome Gambling!"

func RootRoute (c echo.Context) error {
	return c.HTML(http.StatusOK, RootGreeting)
}
