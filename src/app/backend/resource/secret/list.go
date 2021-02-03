package secret

import (
	"context"
	"log"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
)

type SecretSpec interface {
	GetName() string
	GetType() v1.SecretType
	GetNamespace() string
	GetData() map[string][]byte
}

type ImagePullSecretSpec struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Data      []byte `json:"data"`
}

func (spec *ImagePullSecretSpec) GetName() string {
	return spec.Name
}

func (spec *ImagePullSecretSpec) GetType() v1.SecretType {
	return v1.SecretTypeDockercfg
}

func (spec *ImagePullSecretSpec) GetNamespace() string {
	return spec.Namespace
}

func (spec *ImagePullSecretSpec) GetData() map[string][]byte {
	return map[string][]byte{v1.DockerConfigKey: spec.Data}
}

type Secret struct {
	ObjectMeta api.ObjectMeta `json:"objectMeta"`
	TypeMeta   api.TypeMeta   `json:"typeMeta"`
	Type       v1.SecretType  `json:"type"`
}

type SecretList struct {
	api.ListMeta `json:"listMeta"`
	Secrets      []Secret `json:"secrets"`
	Errors       []error  `json:"errors"`
}

func GetSecretList(kubernetes kubernetes.Interface, namespace *common.NamespaceQuery,
	dsQuery *dataselect.DataSelectQuery) (*SecretList, error) {
	log.Printf("Getting list of secrets in %s namespace\n", namespace)
	secretList, err := kubernetes.CoreV1().Secrets(namespace.ToRequestParam()).List(context.TODO(), api.ListEverything)

	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	return ToSecretList(secretList.Items, nonCriticalErrors, dsQuery), nil
}

func CreateSecret(kubernetes kubernetes.Interface, spec SecretSpec) (*Secret, error) {
	namespace := spec.GetNamespace()
	secret := &v1.Secret{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      spec.GetName(),
			Namespace: namespace,
		},
		Type: spec.GetType(),
		Data: spec.GetData(),
	}
	_, err := kubernetes.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metaV1.CreateOptions{})
	result := toSecret(secret)
	return &result, err
}

func toSecret(secret *v1.Secret) Secret {
	return Secret{
		ObjectMeta: api.NewObjectMeta(secret.ObjectMeta),
		TypeMeta:   api.NewTypeMeta(api.ResourceKindSecret),
		Type:       secret.Type,
	}
}

func ToSecretList(secrets []v1.Secret, nonCriticalErrors []error, dsQuery *dataselect.DataSelectQuery) *SecretList {
	newSecretList := &SecretList{
		ListMeta: api.ListMeta{TotalItems: len(secrets)},
		Secrets:  make([]Secret, 0),
		Errors:   nonCriticalErrors,
	}

	secretCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(secrets), dsQuery)
	secrets = fromCells(secretCells)
	newSecretList.ListMeta = api.ListMeta{TotalItems: filteredTotal}

	for _, secret := range secrets {
		newSecretList.Secrets = append(newSecretList.Secrets, toSecret(&secret))
	}

	return newSecretList
}
