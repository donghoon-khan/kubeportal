package ingress

import (
	"reflect"
	"testing"

	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
	networking "k8s.io/api/networking/v1"
)

func TestIngressList(t *testing.T) {
	cases := []struct {
		ingresses []networking.Ingress
		expected  *IngressList
	}{
		{nil, &IngressList{Items: []Ingress{}}},
	}
	for _, c := range cases {
		actual := toIngressList(c.ingresses, nil, dataselect.NoDataSelect)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("toConfigMapList(%#v) == \n%#v\nexpected \n%#v\n",
				c.ingresses, actual, c.expected)
		}
	}
}
