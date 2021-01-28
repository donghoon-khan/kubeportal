package configmap

import (
	"reflect"
	"testing"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestToConfigMapList(t *testing.T) {
	cases := []struct {
		configMaps []v1.ConfigMap
		expected   *ConfigMapList
	}{
		{nil, &ConfigMapList{Items: []ConfigMap{}}},
		{
			[]v1.ConfigMap{
				{Data: map[string]string{"app": "my-name"}, ObjectMeta: metaV1.ObjectMeta{Name: "foo"}},
			},
			&ConfigMapList{
				ListMeta: api.ListMeta{TotalItems: 1},
				Items: []ConfigMap{{
					TypeMeta:   api.TypeMeta{Kind: "configmap"},
					ObjectMeta: api.ObjectMeta{Name: "foo"},
				}},
			},
		},
	}
	for _, c := range cases {
		actual := toConfigMapList(c.configMaps, nil, dataselect.NoDataSelect)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("toConfigMapList(%#v) == \n%#v\nexpected \n%#v\n",
				c.configMaps, actual, c.expected)
		}
	}
}
