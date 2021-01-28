package persistentvolumeclaim

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
)

func TestGetPersistentVolumeClaimList(t *testing.T) {
	cases := []struct {
		persistentVolumeClaims []v1.PersistentVolumeClaim
		expected               *PersistentVolumeClaimList
	}{
		{
			nil,
			&PersistentVolumeClaimList{
				Items: []PersistentVolumeClaim{},
			},
		},
		{
			[]v1.PersistentVolumeClaim{{
				ObjectMeta: metaV1.ObjectMeta{Name: "foo"},
				Spec:       v1.PersistentVolumeClaimSpec{VolumeName: "my-volume"},
				Status:     v1.PersistentVolumeClaimStatus{Phase: v1.ClaimBound},
			}},
			&PersistentVolumeClaimList{
				ListMeta: api.ListMeta{TotalItems: 1},
				Items: []PersistentVolumeClaim{{
					TypeMeta:   api.TypeMeta{Kind: "persistentvolumeclaim"},
					ObjectMeta: api.ObjectMeta{Name: "foo"},
					Status:     "Bound",
					Volume:     "my-volume",
				}},
			},
		},
	}
	for _, c := range cases {
		actual := toPersistentVolumeClaimList(c.persistentVolumeClaims, nil, dataselect.NoDataSelect)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("getPersistentVolumeClaimList(%#v) == \n%#v\nexpected \n%#v\n",
				c.persistentVolumeClaims, actual, c.expected)
		}
	}
}
