package main

import (
	"github.com/artbakulev/techdb/app/server"
	"github.com/artbakulev/techdb/infrastructure"
	"log"
)

func main() {
	connection, err := infrastructure.InitDatabaseConnection()
	if err != nil {
		log.Fatal("Can not connect to database: ", err)
	}
	s := server.NewServer(":8000", connection)
	log.Fatal(s.ListenAndServe())
}
