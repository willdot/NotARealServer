package main

import (
	"log"
	"net/http"

	"github.com/willdot/NotARealServer/handlers"

	"github.com/gorilla/mux"
)

func main() {

	server := handlers.NewPersistServer()

	router := mux.NewRouter()

	router.HandleFunc("/basic", handlers.BasicWithBody())
	router.HandleFunc("/basicwithbody", handlers.BasicWithBody())
	router.HandleFunc("/save", server.SaveRequest())
	router.HandleFunc("/{request}", server.RetreiveRequest())

	err := http.ListenAndServe(":8081", router)

	if err != nil {
		log.Fatal(err)
	}
}
