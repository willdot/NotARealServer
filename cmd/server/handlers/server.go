package handlers

import "github.com/willdot/NotARealServer/pkg/persistrequests"

// Server allows the user to save, retrieve or remove requests
type Server struct {
	FileReadWriter persistrequests.ReadWriter
	FileRemover    persistrequests.Remover
	HandleRequests persistrequests.HandleRequests
}

// NewServer creates a new Server and adds in dependencies
func NewServer(requestDirectory string) *Server {
	return &Server{
		FileReadWriter: persistrequests.FileReadWriter{},
		FileRemover:    persistrequests.FileRemover{},
		HandleRequests: persistrequests.JSONPersist{
			RequestDirectory: requestDirectory,
		},
	}
}
