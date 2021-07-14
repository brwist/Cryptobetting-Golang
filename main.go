package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// config
var (
	DB_NAME     = ""
	DB_HOST     = ""
	DB_USER     = ""
	DB_PASSWORD = ""
	DB_PORT     = ""
	HTTP_PORT   = ""
)

func loadEnvirontment() {

	// init config
	godotenv.Load()

	// init values
	DB_NAME = os.Getenv("DB_NAME")
	DB_HOST = os.Getenv("DB_HOST")
	DB_USER = os.Getenv("DB_USER")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_PORT = os.Getenv("DB_PORT")
	HTTP_PORT = os.Getenv("HTTP_PORT")
}

// connection instance
func createConnection() *gorm.DB {

	// data source
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	// db connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database")
	}

	return db
}

// entities
const (

	// valid fixture, operator can start settlement
	ValidFixture = 0

	//  TBD, reserved for errors in fixtures
	OthersFixture = -1
)

type (
	Fixture struct {
		FixtureID      int64     `gorm:"primaryKey;AUTO_INCREMENT"`
		StartTime      time.Time `gorm:"index;not null"`
		MarketEndTime  time.Time
		EndTime        time.Time `gorm:"index;not null"`
		EndPrice       float64
		FixtureCreated bool
		FixtureEnded   time.Time
		Price          float64
		ExpiryTime     time.Time
		EndFixture     time.Time
		Status         int
	}
)

// schema migration
func migrateSchema() {

	db := createConnection()

	// db logs
	db.Debug()

	// auto migration
	db.AutoMigrate(&Fixture{})
	log.Println("Successfully migrated!")

}

// start scheduler
func startScheduler() {
	sched := gocron.NewScheduler(time.UTC)

	sched.Every(1).Minutes().Do(func() {

		timeNow := time.Now()

		if timeNow.Minute()%15 == 0 {

			// create record
			db := createConnection()
			db.Create(&Fixture{
				FixtureID:      0,
				StartTime:      timeNow.Add(time.Duration(-15) * time.Minute),
				MarketEndTime:  timeNow.Add(time.Duration(-5) * time.Minute),
				EndTime:        timeNow,
				EndPrice:       0,
				FixtureCreated: false,
				FixtureEnded:   time.Time{},
				Price:          0,
				ExpiryTime:     timeNow,
				EndFixture:     time.Time{},
				Status:         OthersFixture,
			})
			log.Println("Fixture created!")
		}

		fmt.Println("Running task at " + timeNow.String())
	})

	sched.StartAsync()
}

// fixture api
type (
	CreateFixtureItem struct {
		Id            int64
		StartTime     string
		MarketEndTime string
		EndTime       string
	}
	CreateFixtureReq struct {
		Timestamp string
		Seq       string
		Fixture   CreateFixtureItem
	}
	CreateFixtureRes struct {
		Timestamp string
		Seq       string
		Status    int
		Message   string
	}
	EndFixtureItem struct {
		Id     int64
		Price  string
		Status int
	}
	EndFixtureReq struct {
		Timestamp string
		Seq       string
		Fixture   EndFixtureItem
	}
	EndFixtureRes struct {
		Timestamp string
		Seq       string
		Status    int
		Message   string
	}
)

func createFixture(req *CreateFixtureReq) CreateFixtureRes {
	var res CreateFixtureRes
	var fixture Fixture
	db := createConnection()

	db.First(&fixture, 1)
	db.First(&fixture, "fixture_id = ?", req.Fixture.Id)

	db.Model(&fixture).Update("EndTime", req.Fixture.EndTime)
	db.Model(&fixture).Update("MarketEndTime", req.Fixture.MarketEndTime)
	db.Model(&fixture).Update("StartTime", req.Fixture.StartTime)
	db.Model(&fixture).Update("Status", OthersFixture)

	res.Timestamp = time.Now().String()
	res.Seq = req.Seq
	res.Message = "Success"
	res.Status = OthersFixture

	return res
}

func endFixture(req *EndFixtureReq) EndFixtureRes {
	var res EndFixtureRes
	var fixture Fixture
	db := createConnection()

	db.First(&fixture, 1)
	db.First(&fixture, "fixture_id = ?", req.Fixture.Id)

	db.Model(&fixture).Update("Price", req.Fixture.Price)
	db.Model(&fixture).Update("Status", ValidFixture)

	res.Timestamp = time.Now().String()
	res.Seq = req.Seq
	res.Message = "Success"
	res.Status = ValidFixture

	return res
}

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

func serveWebsocket(c echo.Context) error {
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

// run server
func runServer() {

	ech := echo.New()

	ech.Use(middleware.Logger())
	ech.Use(middleware.Recover())
	ech.Use(middleware.CORS())

	ech.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Welcome Gambling!")
	})
	ech.POST("/api/CreateFixture", func(ctx echo.Context) error {
		req := &CreateFixtureReq{}
		err := ctx.Bind(req)
		if err != nil {
			return err
		}
		res := createFixture(req)
		return ctx.JSON(http.StatusOK, res)
	})
	ech.POST("/api/EndFixture", func(ctx echo.Context) error {
		req := &EndFixtureReq{}
		err := ctx.Bind(req)
		if err != nil {
			return err
		}
		res := endFixture(req)
		return ctx.JSON(http.StatusOK, res)
	})
	ech.GET("/ws/Fixtures", serveWebsocket)
	ech.Logger.Fatal(ech.Start(":" + HTTP_PORT))

}

func main() {

	// env file
	loadEnvirontment()

	// create schema
	migrateSchema()

	// start scheduler
	startScheduler()

	// http server
	runServer()

}
