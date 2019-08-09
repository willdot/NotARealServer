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

	server := handlers.NewServer(requestFileDirectory)

	router := mux.NewRouter()

	router.HandleFunc("/basic", handlers.BasicHandler())
	router.HandleFunc("/basicwithbody", handlers.BasicWithBodyHandler())
	router.HandleFunc("/save", server.SaveRequestHandler())
	router.HandleFunc("/remove", server.RemoveRequestHandler())
	router.HandleFunc("/removeall", server.RemoveAllRequestsHandler())
	router.HandleFunc("/{RequestRoute}", server.RetreiveRequestHandler())

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func validateRequestDirectory(dir *string) {
	if strings.HasSuffix(*dir, "/") == false {
		*dir += "/"
	}
}
