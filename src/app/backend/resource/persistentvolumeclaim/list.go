package persistentvolumeclaim

import (
	"log"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
)

type PersistentVolumeClaimList struct {
	ListMeta api.ListMeta            `json:"listMeta"`
	Items    []PersistentVolumeClaim `json:"items"`
	Errors   []error                 `json:"errors" swaggertype:"array,string"`
}

type PersistentVolumeClaim struct {
	ObjectMeta   api.ObjectMeta                  `json:"objectMeta"`
	TypeMeta     api.TypeMeta                    `json:"typeMeta"`
	Status       string                          `json:"status"`
	Volume       string                          `json:"volume"`
	Capacity     v1.ResourceList                 `json:"capacity"`
	AccessModes  []v1.PersistentVolumeAccessMode `json:"accessModes"`
	StorageClass *string                         `json:"storageClass"`
}

func GetPersistentVolumeClaimList(client kubernetes.Interface, nsQuery *common.NamespaceQuery,
	dsQuery *dataselect.DataSelectQuery) (*PersistentVolumeClaimList, error) {

	log.Print("Getting list persistent volumes claims")
	channels := &common.ResourceChannels{
		PersistentVolumeClaimList: common.GetPersistentVolumeClaimListChannel(client, nsQuery, 1),
	}

	return GetPersistentVolumeClaimListFromChannels(channels, nsQuery, dsQuery)
}

func GetPersistentVolumeClaimListFromChannels(channels *common.ResourceChannels, nsQuery *common.NamespaceQuery,
	dsQuery *dataselect.DataSelectQuery) (*PersistentVolumeClaimList, error) {

	persistentVolumeClaims := <-channels.PersistentVolumeClaimList.List
	err := <-channels.PersistentVolumeClaimList.Error
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	return toPersistentVolumeClaimList(persistentVolumeClaims.Items, nonCriticalErrors, dsQuery), nil
}

func toPersistentVolumeClaim(pvc v1.PersistentVolumeClaim) PersistentVolumeClaim {
	return PersistentVolumeClaim{
		ObjectMeta:   api.NewObjectMeta(pvc.ObjectMeta),
		TypeMeta:     api.NewTypeMeta(api.ResourceKindPersistentVolumeClaim),
		Status:       string(pvc.Status.Phase),
		Volume:       pvc.Spec.VolumeName,
		Capacity:     pvc.Status.Capacity,
		AccessModes:  pvc.Spec.AccessModes,
		StorageClass: pvc.Spec.StorageClassName,
	}
}

func toPersistentVolumeClaimList(persistentVolumeClaims []v1.PersistentVolumeClaim, nonCriticalErrors []error,
	dsQuery *dataselect.DataSelectQuery) *PersistentVolumeClaimList {

	result := &PersistentVolumeClaimList{
		Items:    make([]PersistentVolumeClaim, 0),
		ListMeta: api.ListMeta{TotalItems: len(persistentVolumeClaims)},
		Errors:   nonCriticalErrors,
	}

	pvcCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(persistentVolumeClaims), dsQuery)
	persistentVolumeClaims = fromCells(pvcCells)
	result.ListMeta = api.ListMeta{TotalItems: filteredTotal}

	for _, item := range persistentVolumeClaims {
		result.Items = append(result.Items, toPersistentVolumeClaim(item))
	}

	return result
}
