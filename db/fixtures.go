package db

import (
	"time"
)

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

func CreateFixture(req *CreateFixtureReq) CreateFixtureRes {
	var res CreateFixtureRes
	var fixture Fixture
	db := CreateConnection()

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

func EndFixture(req *EndFixtureReq) EndFixtureRes {
	var res EndFixtureRes
	var fixture Fixture
	db := CreateConnection()

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
