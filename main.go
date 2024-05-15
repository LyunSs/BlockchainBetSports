package main

import (
	"BlockchainBetSportGolang/bet"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func main(){
	address := os.Args[1]
	router:= bet.NewRouter(address)

	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"POST","GET"})

	log.Fatal(http.ListenAndServe(":"+address,handlers.CORS(allowedOrigins,allowedMethods)(router)))
}