package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-co-op/gocron"
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
			})
			log.Println("Fixture created!")
		}

		fmt.Println("Running task at " + timeNow.String())
	})

	sched.StartAsync()
}

// fixture logic
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

	res.Timestamp = time.Now().String()
	res.Seq = req.Seq
	res.Message = "Success"
	res.Status = fixture.Status

	return res
}

func endFixture(req *EndFixtureReq) EndFixtureRes {
	var res EndFixtureRes
	var fixture Fixture
	db := createConnection()

	db.First(&fixture, 1)
	db.First(&fixture, "fixture_id = ?", req.Fixture.Id)
	db.Model(&fixture).Update("Price", req.Fixture.Price)
	db.Model(&fixture).Update("Status", req.Fixture.Status)

	res.Timestamp = time.Now().String()
	res.Seq = req.Seq
	res.Message = "Success"
	res.Status = req.Fixture.Status

	return res
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
	ech.GET("/api", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Api v1")
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
