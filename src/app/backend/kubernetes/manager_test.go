package kubernetes_test

import (
	"net/http"
	"testing"

	"github.com/donghoon-khan/kubeportal/src/app/backend/kubernetes"
	"github.com/emicklei/go-restful"
)

func TestNewKubernetesManager(t *testing.T) {
	cases := []struct {
		kubeConfigPath, apiserverHost string
	}{
		{"", "test"},
	}

	for _, c := range cases {
		manager := kubernetes.NewKubernetesManager(c.kubeConfigPath, c.apiserverHost)

		if manager == nil {
			t.Fatalf("NewClientManager(%s, %s): Expected manager not to be nil",
				c.kubeConfigPath, c.apiserverHost)
		}
	}

}

func TestKubernetes(t *testing.T) {
	cases := []struct {
		request *restful.Request
	}{
		{
			&restful.Request{
				Request: &http.Request{
					Header: http.Header(map[string][]string{}),
				},
			},
		},
	}

	for _, c := range cases {
		k8sManager := kubernetes.NewKubernetesManager("", "http://localhost:8080")
		_, err := k8sManager.Kubernetes(c.request)

		if err != nil {
			t.Fatalf("Client(%v): Expected client to be created but error was thrown:"+
				" %s", c.request, err.Error())
		}
	}
}
