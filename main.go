package main

import (
	"github.com/onemanshow79/Cryptobetting/db"
	"github.com/onemanshow79/Cryptobetting/server"
	"github.com/onemanshow79/Cryptobetting/tasks"
)

func main() {
	db.Init()
	tasks.StartScheduler()

	// http server
	server.RunServer()
}
