package fixtures

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/onemanshow79/Cryptobetting/db"
	"github.com/onemanshow79/Cryptobetting/models"
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
