package main

import (
	"github.com/artbakulev/techdb/app/server"
	"github.com/artbakulev/techdb/infrastructure"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	connection, err := infrastructure.InitDatabaseConnection()
	if err != nil {
		log.Fatal("Can not connect to database: ", err)
	}
	s := server.NewServer(":5000", connection)
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	log.Fatal(s.ListenAndServe())
}
