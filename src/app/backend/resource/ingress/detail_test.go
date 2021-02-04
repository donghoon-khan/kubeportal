package ingress

import (
	"reflect"
	"testing"

	networking "k8s.io/api/networking/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
)

func TestIngressDetail(t *testing.T) {

	ingressClassName := "kubernetes.io/ingress.class"

	cases := []struct {
		ingress  *networking.Ingress
		expected *IngressDetail
	}{
		{
			&networking.Ingress{
				Spec: networking.IngressSpec{
					IngressClassName: &ingressClassName,
				},
				ObjectMeta: metaV1.ObjectMeta{Name: "foo"},
			},
			&IngressDetail{
				Ingress: Ingress{
					TypeMeta:   api.TypeMeta{Kind: "ingress"},
					ObjectMeta: api.ObjectMeta{Name: "foo"},
					Endpoints:  []common.Endpoint{},
				},
				Spec:   networking.IngressSpec{IngressClassName: &ingressClassName},
				Status: networking.IngressStatus{},
			},
		},
	}
	for _, c := range cases {
		actual := getIngressDetail(c.ingress)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("getIngress(%#v) == \n%#v\nexpected \n%#v\n", c.ingress, actual, c.expected)
		}
	}
}
