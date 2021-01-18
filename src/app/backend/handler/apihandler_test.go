package handler_test

import (
	"testing"

	"github.com/kubernetes/dashboard/src/app/backend/client"

	"github.com/donghoon-khan/kubeportal/src/app/backend/handler"
)

func TestCreateHttpApiHandler(t *testing.T) {
	cManager := client.NewClientManager("", "http://localhost:8080")
	_, err := handler.CreateHttpApiHandler(cManager)
	if err != nil {
		t.Fatal("CreateHttpApiHandler() cannot create HTTP API handler")
	}
}
