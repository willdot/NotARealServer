package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	err := http.ListenAndServe(":8080", router)

	if err != nil {
		log.Fatal(err)
	}
}
