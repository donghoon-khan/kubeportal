package persistentvolumeclaim

import (
	"context"
	"log"
	"strings"

	api "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
)

type PersistentVolumeClaimCell api.PersistentVolumeClaim

func GetPodPersistentVolumeClaims(kubernetes kubernetes.Interface, namespace string, podName string,
	dsQuery *dataselect.DataSelectQuery) (*PersistentVolumeClaimList, error) {

	pod, err := kubernetes.CoreV1().Pods(namespace).Get(context.TODO(), podName, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	claimNames := make([]string, 0)
	if pod.Spec.Volumes != nil && len(pod.Spec.Volumes) > 0 {
		for _, v := range pod.Spec.Volumes {
			persistentVolumeClaim := v.PersistentVolumeClaim
			if persistentVolumeClaim != nil {
				claimNames = append(claimNames, persistentVolumeClaim.ClaimName)
			}
		}
	}

	if len(claimNames) > 0 {
		channels := &common.ResourceChannels{
			PersistentVolumeClaimList: common.GetPersistentVolumeClaimListChannel(
				kubernetes, common.NewSameNamespaceQuery(namespace), 1),
		}

		persistentVolumeClaimList := <-channels.PersistentVolumeClaimList.List
		err = <-channels.PersistentVolumeClaimList.Error
		nonCriticalErrors, criticalError := errors.HandleError(err)
		if criticalError != nil {
			return nil, criticalError
		}

		podPersistentVolumeClaims := make([]api.PersistentVolumeClaim, 0)
		for _, pvc := range persistentVolumeClaimList.Items {
			for _, claimName := range claimNames {
				if strings.Compare(claimName, pvc.Name) == 0 {
					podPersistentVolumeClaims = append(podPersistentVolumeClaims, pvc)
					break
				}
			}
		}

		log.Printf("Found %d persistentvolumeclaims related to %s pod",
			len(podPersistentVolumeClaims), podName)
		return toPersistentVolumeClaimList(podPersistentVolumeClaims,
			nonCriticalErrors, dsQuery), nil
	}

	log.Printf("No persistentvolumeclaims found related to %s pod", podName)
	return &PersistentVolumeClaimList{}, nil
}

func (pvcCell PersistentVolumeClaimCell) GetProperty(name dataselect.PropertyName) dataselect.ComparableValue {
	switch name {
	case dataselect.NameProperty:
		return dataselect.StdComparableString(pvcCell.ObjectMeta.Name)
	case dataselect.CreationTimestampProperty:
		return dataselect.StdComparableTime(pvcCell.ObjectMeta.CreationTimestamp.Time)
	case dataselect.NamespaceProperty:
		return dataselect.StdComparableString(pvcCell.ObjectMeta.Namespace)
	default:
		return nil
	}
}

func toCells(std []api.PersistentVolumeClaim) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = PersistentVolumeClaimCell(std[i])
	}
	return cells
}

func fromCells(cells []dataselect.DataCell) []api.PersistentVolumeClaim {
	std := make([]api.PersistentVolumeClaim, len(cells))
	for i := range std {
		std[i] = api.PersistentVolumeClaim(cells[i].(PersistentVolumeClaimCell))
	}
	return std
}
