package main

import (
	"log"
	"net/http"

	"github.com/willdot/NotARealServer/handlers"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/basic", handlers.BasicWithBody())
	router.HandleFunc("/basicwithbody", handlers.BasicWithBody())
	router.HandleFunc("/save", handlers.SaveRequest())
	router.HandleFunc("/{request}", handlers.RetreiveRequest())

	err := http.ListenAndServe(":8081", router)

	if err != nil {
		log.Fatal(err)
	}
}
