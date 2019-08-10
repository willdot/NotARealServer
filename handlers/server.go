package handlers

import "github.com/willdot/NotARealServer/persistrequests"

// Server allows the user to save, retrieve or remove requests
type Server struct {
	//FileWriter     persistrequests.Writer
	//FileReader     persistrequests.Reader
	FileReadWriter persistrequests.ReadWriter
	FileRemover    persistrequests.Remover
	HandleRequests persistrequests.HandleRequests
}

// NewServer creates a new Server and adds in dependencies
func NewServer(requestDirectory string) Server {
	return Server{
		//FileWriter:  persistrequests.FileWriter{},
		//FileReader:  persistrequests.FileReader{},
		FileReadWriter: persistrequests.FileReadWriter{},
		FileRemover:    persistrequests.FileRemover{},
		HandleRequests: persistrequests.JSONPersist{
			RequestDirectory: requestDirectory,
		},
	}
}
