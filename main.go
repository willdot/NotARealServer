package main

import (
	"log"
	"net/http"
	"os"

	"github.com/willdot/NotARealServer/handlers"

	"github.com/gorilla/mux"
)

func main() {

	var port string

	if port = os.Getenv("PORT"); port == "" {
		port = "8080"
	}

	server := handlers.NewPersistServer()

	router := mux.NewRouter()

	router.HandleFunc("/basic", handlers.BasicWithBody())
	router.HandleFunc("/basicwithbody", handlers.BasicWithBody())
	router.HandleFunc("/save", server.SaveRequestHandler())
	router.HandleFunc("/{RequestRoute}", server.RetreiveRequestHandler())

	err := http.ListenAndServe(":"+port, router)

	if err != nil {
		log.Fatal(err)
	}
}
