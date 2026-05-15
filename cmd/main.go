package main

import (
	"log"

	"github.com/Mikhail-tal63/viltrum_empier/cmd/api"
	"github.com/Mikhail-tal63/viltrum_empier/config"
	"github.com/Mikhail-tal63/viltrum_empier/database"
)

func main() {

	mongoDB, err := database.NewMongoStorage(config.ENVs.MongoURI, config.ENVs.DBName)
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(":"+config.ENVs.Port, mongoDB)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
