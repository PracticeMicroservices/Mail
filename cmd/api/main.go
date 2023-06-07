package main

import (
	"log"
)

func main() {
	app := NewApp()

	log.Println("Starting Mail service on port 80")

	//start server
	app.serve()
}
