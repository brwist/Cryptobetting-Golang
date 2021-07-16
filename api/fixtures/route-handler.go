package fixtures

import "github.com/labstack/echo/v4"

func HandleRoutes(ech *echo.Echo) {
	ech.POST("/api/CreateFixture", CreateFixture)
	ech.POST("/api/EndFixture", EndFixture)
}
