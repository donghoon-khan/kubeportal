package ingress

import (
	"context"
	"log"

	networking "k8s.io/api/networking/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

type IngressDetail struct {
	Ingress `json:",inline"`
	Spec    networking.IngressSpec   `json:"spec"`
	Status  networking.IngressStatus `json:"status"`
	Errors  []error                  `json:"errors"`
}

func GetIngressDetail(kubernetes kubernetes.Interface, namespace, name string) (*IngressDetail, error) {
	log.Printf("Getting details of %s ingress in %s namespace", name, namespace)

	rawIngress, err := kubernetes.NetworkingV1().Ingresses(namespace).Get(context.TODO(), name, metaV1.GetOptions{})

	if err != nil {
		return nil, err
	}

	return getIngressDetail(rawIngress), nil
}

func getIngressDetail(i *networking.Ingress) *IngressDetail {
	return &IngressDetail{
		Ingress: toIngress(i),
		Spec:    i.Spec,
		Status:  i.Status,
	}
}
