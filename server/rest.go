package server

import (
	"github.com/allgoodworks/Cryptobetting-Golang/api"
	"github.com/allgoodworks/Cryptobetting-Golang/api/fixtures"
	"github.com/allgoodworks/Cryptobetting-Golang/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// run server
func RunServer() {
	ech := echo.New()

	ech.Use(middleware.Logger())
	ech.Use(middleware.Recover())
	ech.Use(middleware.CORS())

	api.HandleRoutes(ech)
	fixtures.HandleRoutes(ech)

	HandleWebsocketRoutes(ech)

	// public assets
	ech.File("/Fixtures-Test", "public/fixtures-test.htm")

	ech.Logger.Fatal(ech.Start(":" + db.HTTP_PORT))
}
