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
	}
)

// create connection
func createConnection() *gorm.DB {

	// connection string
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")

	// data source
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	// db connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database")
	}

	return db
}

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
			db.Create(&Fixture{StartTime: timeNow.Add(time.Duration(-15) * time.Minute), EndTime: timeNow, MarketEndTime: timeNow.Add(time.Duration(-5) * time.Minute)})
			log.Println("Fixture created!")
		}

		fmt.Println("Running task at " + timeNow.String())
	})

	sched.StartAsync()
}

// run server
func runServer() {

	ech := echo.New()
	ech.Use(middleware.Logger())
	ech.Use(middleware.Recover())
	ech.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Welcome Gambling API V1!")
	})
	ech.POST("/CreateFixture", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Create Fixture")
	})
	ech.POST("/EndFixture", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "End Fixture")
	})
	httpPort := os.Getenv("HTTP_PORT")
	ech.Logger.Fatal(ech.Start(":" + httpPort))

}

func main() {

	// env file
	godotenv.Load()

	// create schema
	migrateSchema()

	// start scheduler
	startScheduler()

	// http server
	runServer()

}
