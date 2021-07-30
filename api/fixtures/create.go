package fixtures

import (
	"net/http"

	"github.com/allgoodworks/Cryptobetting-Golang/db"
	"github.com/allgoodworks/Cryptobetting-Golang/models"
	"github.com/labstack/echo/v4"
)

func CreateFixture (ctx echo.Context) error {
	req := &models.CreateFixtureReq{}
	err := ctx.Bind(req)
	if err != nil {
		return err
	}
	db := db.CreateConnection()

	res := models.CreateFixture(req, db)
	return ctx.JSON(http.StatusOK, res)
}
