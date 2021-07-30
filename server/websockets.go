package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/onemanshow79/Cryptobetting/db"
	"github.com/onemanshow79/Cryptobetting/models"
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
	upgrader = websocket.Upgrader{
		ReadBufferSize:  102400,
		WriteBufferSize: 102400,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
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

		// wating for 1seconds
		time.Sleep(time.Second)

		// Write
		db := db.CreateConnection()
		fixtures := models.GetFixtures(db)
		var results []WsFixtureItem
		for _, fixture := range fixtures {
			item := WsFixtureItem{}
			item.FixtureId = fixture.FixtureID
			results = append(results, item)
		}
		response := &WsFixtureRes{
			Timestamp: time.Now().String(),
			Price:     0,
			Type:      OnlyPriceType,
			Fixtures:  results,
		}
		data, _ := json.Marshal(response)
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
