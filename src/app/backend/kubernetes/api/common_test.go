package api_test

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/authorization/v1"

	"github.com/donghoon-khan/kubeportal/src/app/backend/kubernetes/api"
)

func TestToSelfSubjectAccessReview(t *testing.T) {
	namespace := "test-namespace"
	name := "test-name"
	resourceName := "deployment"
	verb := "GET"
	expected := &v1.SelfSubjectAccessReview{
		Spec: v1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &v1.ResourceAttributes{
				Namespace: namespace,
				Name:      name,
				Resource:  "deployments",
				Verb:      "get",
			},
		},
	}

	got := api.ToSelfSubjectAccessReview(namespace, name, resourceName, verb)
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Expected to get %+v but got %+v", expected, got)
	}
}
