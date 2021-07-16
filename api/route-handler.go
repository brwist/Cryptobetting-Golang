package api

import "github.com/labstack/echo/v4"

func HandleRoutes(ech *echo.Echo) {
	ech.GET("/", RootRoute)
}
