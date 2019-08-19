package handlers

import (
	"reflect"
	"testing"

	"github.com/willdot/NotARealServer/pkg/persistrequests"
)

func TestNewPersistServer(t *testing.T) {

	got := NewServer("")

	want := &Server{
		HandleRequests: persistrequests.JSONPersist{},
		FileReadWriter: persistrequests.FileReadWriter{},
		FileRemover:    persistrequests.FileRemover{},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}
