package persistentvolumeclaim

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
)

func TestGetPodPersistentVolumeClaims(t *testing.T) {
	cases := []struct {
		pod                       *v1.Pod
		name                      string
		namespace                 string
		persistentVolumeClaimList *v1.PersistentVolumeClaimList
		expected                  *PersistentVolumeClaimList
	}{
		{
			pod: &v1.Pod{
				ObjectMeta: metaV1.ObjectMeta{
					Name: "test-pod", Namespace: "test-namespace", Labels: map[string]string{"app": "test"},
				},
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{{
						Name: "vol-1",
						VolumeSource: v1.VolumeSource{
							PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
								ClaimName: "pvc-1",
							},
						},
					}},
				},
			},
			name:      "test-pod",
			namespace: "test-namespace",
			persistentVolumeClaimList: &v1.PersistentVolumeClaimList{Items: []v1.PersistentVolumeClaim{
				{
					ObjectMeta: metaV1.ObjectMeta{
						Name: "pvc-1", Namespace: "test-namespace", Labels: map[string]string{"app": "test"},
					},
				},
			}},
			expected: &PersistentVolumeClaimList{
				ListMeta: api.ListMeta{TotalItems: 1},
				Items: []PersistentVolumeClaim{{
					TypeMeta: api.TypeMeta{Kind: api.ResourceKindPersistentVolumeClaim},
					ObjectMeta: api.ObjectMeta{Name: "pvc-1", Namespace: "test-namespace",
						Labels: map[string]string{"app": "test"}},
				}},
				Errors: []error{},
			},
		},
	}

	for _, c := range cases {

		fakeClient := fake.NewSimpleClientset(c.persistentVolumeClaimList, c.pod)

		actual, _ := GetPodPersistentVolumeClaims(fakeClient, c.namespace, c.name, dataselect.NoDataSelect)

		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("GetPodPersistentVolumeClaims(client, %#v, %#v) == \ngot: %#v, \nexpected %#v",
				c.name, c.namespace, actual, c.expected)
		}
	}
}
