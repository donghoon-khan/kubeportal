package serviceaccount

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/secret"
)

func GetServiceAccountImagePullSecrets(kubernetes kubernetes.Interface, namespace,
	name string, dsQuery *dataselect.DataSelectQuery) (*secret.SecretList, error) {
	imagePullSecretList := secret.SecretList{
		Secrets: []secret.Secret{},
	}

	serviceAccount, err := kubernetes.CoreV1().ServiceAccounts(namespace).Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return &imagePullSecretList, err
	}

	if serviceAccount.ImagePullSecrets == nil {
		return &imagePullSecretList, nil
	}

	channels := &common.ResourceChannels{
		SecretList: common.GetSecretListChannel(kubernetes, common.NewSameNamespaceQuery(namespace), 1),
	}

	apiSecretList := <-channels.SecretList.List
	if err := <-channels.SecretList.Error; err != nil {
		return &imagePullSecretList, err
	}

	imagePullSecretsMap := map[string]struct{}{}
	for _, ips := range serviceAccount.ImagePullSecrets {
		imagePullSecretsMap[ips.Name] = struct{}{}
	}

	var rawImagePullSecretList []v1.Secret
	for _, apiSecret := range apiSecretList.Items {
		if _, ok := imagePullSecretsMap[apiSecret.Name]; ok {
			rawImagePullSecretList = append(rawImagePullSecretList, apiSecret)
		}
	}

	return secret.ToSecretList(rawImagePullSecretList, []error{}, dsQuery), nil
}

func GetServiceAccountSecrets(kubernetes kubernetes.Interface, namespace,
	name string, dsQuery *dataselect.DataSelectQuery) (*secret.SecretList, error) {
	secretList := secret.SecretList{
		Secrets: []secret.Secret{},
	}

	serviceAccount, err := kubernetes.CoreV1().ServiceAccounts(namespace).Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return &secretList, err
	}

	if serviceAccount.Secrets == nil {
		return &secretList, nil
	}

	channels := &common.ResourceChannels{
		SecretList: common.GetSecretListChannel(kubernetes, common.NewSameNamespaceQuery(namespace), 1),
	}

	apiSecretList := <-channels.SecretList.List
	if err := <-channels.SecretList.Error; err != nil {
		return &secretList, err
	}

	secretsMap := map[string]v1.ObjectReference{}
	for _, s := range serviceAccount.Secrets {
		secretsMap[s.Name] = s
	}

	var rawSecretList []v1.Secret
	for _, apiSecret := range apiSecretList.Items {
		if _, ok := secretsMap[apiSecret.Name]; ok {
			rawSecretList = append(rawSecretList, apiSecret)
		}
	}

	return secret.ToSecretList(rawSecretList, []error{}, dsQuery), nil
}
