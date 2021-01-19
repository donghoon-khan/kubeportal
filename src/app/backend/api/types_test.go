package api_test

import (
	"testing"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestIsSelectorMatching(t *testing.T) {
	cases := []struct {
		serviceSelector, replicationControllerSelector map[string]string
		expected                                       bool
	}{
		{nil, nil, false},
		{map[string]string{}, map[string]string{}, false},
		{map[string]string{"app": "kube-portal"}, map[string]string{}, false},
		{map[string]string{"app": "kube-portal"}, map[string]string{"app": "portal"}, false},
		{map[string]string{"app": "kube-portal", "version": "1.1"},
			map[string]string{"app": "kube-portal"}, false},
		{map[string]string{"app": "kube-portal", "version": "1.1"},
			map[string]string{"app": "kube-portal", "version": "1.1"}, true},
		{map[string]string{"app": "kube-portal"},
			map[string]string{"app": "kube-portal", "version": "1.1"}, true},
	}
	for _, c := range cases {
		actual := api.IsSelectorMatching(c.serviceSelector, c.replicationControllerSelector)
		if actual != c.expected {
			t.Errorf("isSelectorMatching(%+vl %+v) == %+v, expected %+v",
				c.serviceSelector, c.replicationControllerSelector, actual, c.expected)
		}
	}
}

func TestIsLabelSelectorMatching(t *testing.T) {
	cases := []struct {
		serviceSelector   map[string]string
		daemonSetSelector *metaV1.LabelSelector
		expected          bool
	}{
		{nil, nil, false},
		{nil, &metaV1.LabelSelector{MatchLabels: map[string]string{}}, false},
		{map[string]string{}, nil, false},
		{map[string]string{}, &metaV1.LabelSelector{MatchLabels: map[string]string{}},
			false},
		{map[string]string{"app": "my-name"},
			&metaV1.LabelSelector{MatchLabels: map[string]string{}},
			false},
		{map[string]string{"app": "my-name", "version": "2"},
			&metaV1.LabelSelector{MatchLabels: map[string]string{"app": "my-name", "version": "1.1"}},
			false},
		{map[string]string{"app": "my-name", "env": "prod"},
			&metaV1.LabelSelector{MatchLabels: map[string]string{"app": "my-name", "version": "1.1"}},
			false},
		{map[string]string{"app": "my-name"},
			&metaV1.LabelSelector{MatchLabels: map[string]string{"app": "my-name"}},
			true},
		{map[string]string{"app": "my-name", "version": "1.1"},
			&metaV1.LabelSelector{MatchLabels: map[string]string{"app": "my-name", "version": "1.1"}},
			true},
	}
	for _, c := range cases {
		actual := api.IsLabelSelectorMatching(c.serviceSelector, c.daemonSetSelector)
		if actual != c.expected {
			t.Errorf("isLabelSelectorMatching(%+v, %+v) == %+v, expected %+v",
				c.serviceSelector, c.daemonSetSelector, actual, c.expected)
		}
	}
}
