package persistentvolumeclaim

import (
	"context"
	"log"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PersistentVolumeClaimDetail struct {
	PersistentVolumeClaim `json:",inline"`
}

func GetPersistentVolumeClaimDetail(kubernetes kubernetes.Interface, namespace string, name string) (*PersistentVolumeClaimDetail, error) {
	log.Printf("Getting details of %s persistent volume claim", name)

	pvc, err := kubernetes.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return getPersistentVolumeClaimDetail(*pvc), nil
}

func getPersistentVolumeClaimDetail(pvc v1.PersistentVolumeClaim) *PersistentVolumeClaimDetail {
	return &PersistentVolumeClaimDetail{
		PersistentVolumeClaim: toPersistentVolumeClaim(pvc),
	}
}
