package pod

import (
	"encoding/base64"
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	metricApi "github.com/donghoon-khan/kubeportal/src/app/backend/integration/metric/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/controller"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/persistentvolumeclaim"
)

func TestGetPodDetail(t *testing.T) {
	cases := []struct {
		pod      *v1.PodList
		expected *PodDetail
	}{
		{
			pod: &v1.PodList{Items: []v1.Pod{{
				ObjectMeta: metaV1.ObjectMeta{
					Name: "test-pod", Namespace: "test-namespace",
					Labels: map[string]string{"app": "test"},
				}}}},
			expected: &PodDetail{
				TypeMeta: api.TypeMeta{Kind: api.ResourceKindPod},
				PodPhase: string(v1.PodUnknown),
				ObjectMeta: api.ObjectMeta{
					Name:      "test-pod",
					Namespace: "test-namespace",
					Labels:    map[string]string{"app": "test"},
				},
				Controller:     &controller.ResourceOwner{},
				Containers:     []Container{},
				InitContainers: []Container{},
				EventList: common.EventList{
					Events: []common.Event{},
					Errors: []error{},
				},
				Metrics:                   []metricApi.Metric{},
				PersistentvolumeclaimList: persistentvolumeclaim.PersistentVolumeClaimList{},
				Errors:                    []error{},
			},
		},
	}

	for _, c := range cases {
		fakeClient := fake.NewSimpleClientset(c.pod)

		dataselect.DefaultDataSelectWithMetrics.MetricQuery = dataselect.NoMetrics
		actual, err := GetPodDetail(fakeClient, nil, "test-namespace", "test-pod")

		if err != nil {
			t.Errorf("GetPodDetail(%#v) == \ngot err %#v", c.pod, err)
		}
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("GetPodDetail(%#v) == \ngot %#v, \nexpected %#v", c.pod, actual, c.expected)
		}
	}
}

func TestEvalValueFrom(t *testing.T) {
	cases := []struct {
		src        *v1.EnvVarSource
		container  *v1.Container
		pod        *v1.Pod
		configMaps *v1.ConfigMapList
		secrets    *v1.SecretList
		expected   string
	}{
		{
			src: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: "secret-env",
					},
					Key: "username",
				},
			},
			container:  nil,
			pod:        nil,
			configMaps: nil,
			secrets: &v1.SecretList{
				Items: []v1.Secret{
					{
						ObjectMeta: metaV1.ObjectMeta{
							Name: "secret-env",
						},
						Data: map[string][]byte{
							"username": []byte("top-secret"),
						},
					},
				},
			},
			expected: base64.StdEncoding.EncodeToString([]byte("top-secret")),
		},
		{
			src: &v1.EnvVarSource{
				ConfigMapKeyRef: &v1.ConfigMapKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: "config-map-env",
					},
					Key: "username",
				},
			},
			container: nil,
			pod:       nil,
			configMaps: &v1.ConfigMapList{
				Items: []v1.ConfigMap{
					{
						ObjectMeta: metaV1.ObjectMeta{
							Name: "config-map-env",
						},
						Data: map[string]string{
							"username": "joey",
						},
					},
				},
			},
			secrets:  nil,
			expected: "joey",
		},
	}

	for _, c := range cases {
		dataselect.DefaultDataSelectWithMetrics.MetricQuery = dataselect.NoMetrics
		actual := evalValueFrom(c.src, c.container, c.pod, c.configMaps, c.secrets)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("GetPodDetail(%#v, %#v, %#v, %#v, %#v) == \ngot %#v, \nexpected %#v",
				c.src, c.container, c.pod, c.configMaps, c.secrets, actual, c.expected)
		}
	}
}

func TestEvalEnvFrom(t *testing.T) {
	cases := []struct {
		container  v1.Container
		configMaps *v1.ConfigMapList
		secrets    *v1.SecretList
		expected   []EnvVar
	}{
		{
			container: v1.Container{
				Name:  "echoserver",
				Image: "k8s.gcr.io/echoserver",
				EnvFrom: []v1.EnvFromSource{
					{
						SecretRef: &v1.SecretEnvSource{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "secret-env",
							},
						},
					}, {
						Prefix: "test_",
						ConfigMapRef: &v1.ConfigMapEnvSource{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "config-map-env",
							},
						},
					},
				},
			},
			configMaps: &v1.ConfigMapList{
				Items: []v1.ConfigMap{
					{
						ObjectMeta: metaV1.ObjectMeta{
							Name: "config-map-env",
						},
						Data: map[string]string{
							"username": "joey",
						},
					},
				},
			},
			secrets: &v1.SecretList{
				Items: []v1.Secret{
					{
						ObjectMeta: metaV1.ObjectMeta{
							Name: "secret-env",
						},
						Data: map[string][]byte{
							"username": []byte("top-secret"),
						},
					},
				},
			},
			expected: []EnvVar{
				{
					Name:  "username",
					Value: base64.StdEncoding.EncodeToString([]byte("top-secret")),
					ValueFrom: &v1.EnvVarSource{
						SecretKeyRef: &v1.SecretKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "secret-env",
							},
							Key: "username",
						},
					},
				},
				{
					Name:  "test_username",
					Value: "joey",
					ValueFrom: &v1.EnvVarSource{
						ConfigMapKeyRef: &v1.ConfigMapKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "config-map-env",
							},
							Key: "username",
						},
					},
				},
			},
		},
	}

	for _, c := range cases {
		dataselect.DefaultDataSelectWithMetrics.MetricQuery = dataselect.NoMetrics
		actual := evalEnvFrom(c.container, c.configMaps, c.secrets)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("evalEnvFrom(%#v, %#v, %#v) == \ngot %#v, \nexpected %#v",
				c.container, c.configMaps, c.secrets, actual, c.expected)
		}
	}
}
