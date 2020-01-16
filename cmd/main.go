package main

import (
	"github.com/artbakulev/techdb/infrastructure"
	"log"
)

func main() {
	err := infrastructure.InitDatabaseConnection()
	if err != nil {
		log.Fatal("Can not connect to database: ", err)
	}
}
