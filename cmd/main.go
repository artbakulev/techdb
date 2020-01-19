package main

import (
	server2 "github.com/artbakulev/techdb/app/server"
	"github.com/artbakulev/techdb/app/user/delivery/http"
	"github.com/artbakulev/techdb/app/user/repository"
	"github.com/artbakulev/techdb/app/user/usecase"
	"github.com/artbakulev/techdb/infrastructure"
	"github.com/buaazp/fasthttprouter"
	"log"
)

func main() {
	connection, err := infrastructure.InitDatabaseConnection()
	if err != nil {
		log.Fatal("Can not connect to database: ", err)
	}

	userRepo := repository.NewPostgresUserRepository(connection)
	userUsecase := usecase.NewUserUsecase(userRepo)

	router := fasthttprouter.New()

	http.NewUserHandler(router, userUsecase)

	server := server2.NewServer(":8000", router)
	log.Fatal(server.ListenAndServe())
}
