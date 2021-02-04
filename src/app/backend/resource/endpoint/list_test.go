package endpoint

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
)

func TestToEndpointList(t *testing.T) {
	cases := []struct {
		endpoints []v1.Endpoints
		expected  *EndpointList
	}{
		{nil, &EndpointList{Endpoints: []Endpoint{}}},
	}
	for _, c := range cases {
		actual := toEndpointList(c.endpoints)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("toEndpointList(%#v) == \n%#v\nexpected \n%#v\n", c.endpoints, actual, c.expected)
		}
	}
}
