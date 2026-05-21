package main

import (
	"log"

	"github.com/Mikhail-tal63/viltrum_empier/cmd/api"
	"github.com/Mikhail-tal63/viltrum_empier/config"
	"github.com/Mikhail-tal63/viltrum_empier/database"
	"github.com/Mikhail-tal63/viltrum_empier/internal/websocket"
)

func main() {

	mongoDB, err := database.NewMongoStorage(config.ENVs.MongoURI, config.ENVs.DBName)
	if err != nil {
		log.Fatal(err)
	}

	hub := websocket.NewHub()
	go hub.Run()

	server := api.NewAPIServer(":"+config.ENVs.Port, mongoDB, hub)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
