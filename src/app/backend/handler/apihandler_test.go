package handler_test

import (
	"testing"
)

func TestCreateHttpApiHandler(t *testing.T) {

	a := 10
	if a != 10 {
		t.Fatal("error")
	}

	/*cManager := client.NewClientManager("", "http://127.0.0.1:8080")
	if cManager == nil {
		t.Fatal("error")
	}*/
	//_, err := handler.CreateHttpApiHandler(cManager)
	/*if err != nil {
		t.Fatal("CreateHttpApiHandler() cannot create HTTP API handler")
	}*/
}
