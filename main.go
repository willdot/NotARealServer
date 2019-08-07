package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/willdot/NotARealServer/handlers"

	"github.com/gorilla/mux"
)

func main() {

	var port string

	if port = os.Getenv("PORT"); port == "" {
		port = "8080"
	}

	var requestFileDirectory string
	if requestFileDirectory = os.Getenv("REQUESTDIRECTORY"); requestFileDirectory == "" {
		requestFileDirectory = "requests/"
	}

	validateRequestDirectory(&requestFileDirectory)

	server := handlers.NewPersistServer(requestFileDirectory)

	router := mux.NewRouter()

	router.HandleFunc("/basic", handlers.Basic())
	router.HandleFunc("/basicwithbody", handlers.BasicWithBody())
	router.HandleFunc("/save", server.SaveRequestHandler())
	router.HandleFunc("/{RequestRoute}", server.RetreiveRequestHandler())

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func validateRequestDirectory(dir *string) {
	if strings.HasSuffix(*dir, "/") == false {
		*dir += "/"
	}
}
