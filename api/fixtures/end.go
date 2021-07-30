package fixtures

import (
	"net/http"

	"github.com/allgoodworks/Cryptobetting-Golang/db"
	"github.com/allgoodworks/Cryptobetting-Golang/models"
	"github.com/labstack/echo/v4"
)

func EndFixture(ctx echo.Context) error {
	req := &models.EndFixtureReq{}
	err := ctx.Bind(req)
	if err != nil {
		return err
	}
	db := db.CreateConnection()

	res := models.EndFixture(req, db)
	return ctx.JSON(http.StatusOK, res)
}
