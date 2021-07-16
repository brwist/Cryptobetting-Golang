package server

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// fixture socket
const (
	// price only
	OnlyPriceType = 1

	//  price + lines
	PriceLineType = 2
)

type (
	WsFixtureProbability struct {
		Strike float64
		Over   float32
		Under  float32
	}
	WsFixtureItem struct {
		FixtureId     int64
		Probabilities []WsFixtureProbability
	}
	WsFixtureRes struct {
		Timestamp string
		Price     float64
		Type      int
		Fixtures  []WsFixtureItem
	}
)

var (
	upgrader     = websocket.Upgrader{}
	wsFixtureRes = &WsFixtureRes{
		Timestamp: time.Now().String(),
		Price:     0,
		Type:      OnlyPriceType,
		Fixtures:  []WsFixtureItem{},
	}
)

func HandleWebsocketRoutes(ech *echo.Echo) {
	ech.GET("/ws/Fixtures", ServeWebsocket)
}

func ServeWebsocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		//store
		data, _ := json.Marshal(wsFixtureRes)

		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte(data))
		if err != nil {
			c.Logger().Error(err)
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)
	}
}
