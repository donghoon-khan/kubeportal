package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/emicklei/go-restful/v3"

	"github.com/donghoon-khan/kubeportal/src/app/backend/args"
)

/*func TestCreateHttpApiHandler(t *testing.T) {

	kManager := kubernetes.NewKubernetesManager("", "http://localhost:8080")
	aManager := auth.NewAuthManager(kManager, authApi.AuthenticationModes{}, true)

	_, err := CreateHttpApiHandler(kManager, aManager)
	if err != nil {
		t.Fatal("CreateHttpApiHandler() cannot create HTTP API handler")
	}

	cManager := client.NewClientManager("", "http://127.0.0.1:8080")
	if cManager == nil {
		t.Fatal("error")
	}
	_, err := handler.CreateHttpApiHandler(cManager)
	if err != nil {
		t.Fatal("CreateHttpApiHandler() cannot create HTTP API handler")
	}
}*/

func TestMapUrlToResource(t *testing.T) {
	cases := []struct {
		url, expected string
	}{
		{
			"/api/v1/pod",
			"pod",
		},
		{
			"/api/v1/node",
			"node",
		},
	}
	for _, c := range cases {
		actual := mapUrlToResource(c.url)
		if !reflect.DeepEqual(actual, &c.expected) {
			t.Errorf("mapUrlToResource(%#v) returns %#v, expected %#v", c.url, *actual, c.expected)
		}
	}
}

func TestFormatRequestLog(t *testing.T) {
	cases := []struct {
		method      string
		uri         string
		content     map[string]string
		expected    string
		apiLogLevel string
	}{
		{
			"PUT",
			"/api/v1/pod",
			map[string]string{},
			"Incoming HTTP/1.1 PUT /api/v1/pod request",
			"DEFAULT",
		},
		{
			"PUT",
			"/api/v1/pod",
			map[string]string{},
			"",
			"NONE",
		},
		{
			"POST",
			"/api/v1/login",
			map[string]string{"password": "abc123"},
			"Incoming HTTP/1.1 POST /api/v1/login request from : { contents hidden }",
			"DEFAULT",
		},
		{
			"POST",
			"/api/v1/login",
			map[string]string{},
			"",
			"NONE",
		},
		{
			"POST",
			"/api/v1/login",
			map[string]string{"password": "abc123"},
			"Incoming HTTP/1.1 POST /api/v1/login request from : {\"password\":\"abc123\"}",
			"DEBUG",
		},
	}

	for _, c := range cases {
		jsonValue, _ := json.Marshal(c.content)

		req, err := http.NewRequest(c.method, c.uri, bytes.NewReader(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		if err != nil {
			t.Error("Cannot mockup request")
		}

		builder := args.GetHolderBuilder()
		builder.SetApiLogLevel(c.apiLogLevel)

		var restfulRequest restful.Request
		restfulRequest.Request = req

		actual := formatRequestLog(&restfulRequest)
		if !strings.Contains(actual, c.expected) {
			t.Errorf("formatRequestLog(%#v) returns %#v, expected to contain %#v", req, actual, c.expected)
		}
	}
}
