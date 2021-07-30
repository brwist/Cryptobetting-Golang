package tasks

import (
	"fmt"
	"log"
	"time"

	"github.com/allgoodworks/Cryptobetting-Golang/db"
	"github.com/allgoodworks/Cryptobetting-Golang/models"
	"github.com/go-co-op/gocron"
)

// Start scheduler
func StartScheduler() {
	sched := gocron.NewScheduler(time.UTC)

	sched.Every(1).Minutes().Do(func() {

		timeNow := time.Now()

		if timeNow.Minute()%15 == 0 {

			// create record
			dbConn := db.CreateConnection()
			dbConn.Create(&models.Fixture{
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
				Status:         models.OthersFixture,
			})
			log.Println("Fixture created!")
		}

		fmt.Println("Running task at " + timeNow.String())
	})

	sched.StartAsync()
}
