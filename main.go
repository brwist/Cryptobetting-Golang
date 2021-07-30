package main

import (
	"github.com/allgoodworks/Cryptobetting-Golang/db"
	"github.com/allgoodworks/Cryptobetting-Golang/server"
	"github.com/allgoodworks/Cryptobetting-Golang/tasks"
)

func main() {

	// db migrator
	db.Init()

	// sync jobs
	tasks.StartScheduler()

	// http server
	server.RunServer()
}
