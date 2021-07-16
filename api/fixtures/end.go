package fixtures

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/onemanshow79/Cryptobetting/db"
)

func EndFixture (ctx echo.Context) error {
	req := &db.EndFixtureReq{}
	err := ctx.Bind(req)
	if err != nil {
		return err
	}
	res := db.EndFixture(req)
	return ctx.JSON(http.StatusOK, res)
}
