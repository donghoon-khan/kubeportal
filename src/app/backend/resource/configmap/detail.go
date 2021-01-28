package configmap

import (
	"context"
	"log"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ConfigMapDetail struct {
	ConfigMap `json:",inline"`
	Data      map[string]string `json:"data,omitempty"`
}

func GetConfigMapDetail(client kubernetes.Interface, namespace, name string) (*ConfigMapDetail, error) {
	log.Printf("Getting details of %s config map in %s namespace", name, namespace)

	rawConfigMap, err := client.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return getConfigMapDetail(rawConfigMap), nil
}

func getConfigMapDetail(rawConfigMap *v1.ConfigMap) *ConfigMapDetail {
	return &ConfigMapDetail{
		ConfigMap: toConfigMap(rawConfigMap.ObjectMeta),
		Data:      rawConfigMap.Data,
	}
}
