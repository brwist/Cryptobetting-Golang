package fixtures

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/onemanshow79/Cryptobetting/db"
)

func CreateFixture (ctx echo.Context) error {
	req := &db.CreateFixtureReq{}
	err := ctx.Bind(req)
	if err != nil {
		return err
	}
	res := db.CreateFixture(req)
	return ctx.JSON(http.StatusOK, res)
}
