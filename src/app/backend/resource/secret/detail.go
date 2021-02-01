package secret

import (
	"context"
	"log"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type SecretDetail struct {
	Secret `json:",inline"`
	Data   map[string][]byte `json:"data"`
}

func GetSecretDetail(kubernetes kubernetes.Interface, namespace, name string) (*SecretDetail, error) {
	log.Printf("Getting details of %s secret in %s namespace\n", name, namespace)

	rawSecret, err := kubernetes.CoreV1().Secrets(namespace).Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return getSecretDetail(rawSecret), nil
}

func getSecretDetail(rawSecret *v1.Secret) *SecretDetail {
	return &SecretDetail{
		Secret: toSecret(rawSecret),
		Data:   rawSecret.Data,
	}
}
