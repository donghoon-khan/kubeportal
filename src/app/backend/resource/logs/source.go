package logs

import (
	"context"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/controller"
)

func GetLogSources(kubernetes kubernetes.Interface, ns, resourceName, resourceType string) (controller.LogSources, error) {
	if resourceType == "pod" {
		return getLogSourcesFromPod(kubernetes, ns, resourceName)
	}
	return getLogSourcesFromController(kubernetes, ns, resourceName, resourceType)
}

func getLogSourcesFromPod(kubernetes kubernetes.Interface, ns, resourceName string) (controller.LogSources, error) {
	pod, err := kubernetes.CoreV1().Pods(ns).Get(context.TODO(), resourceName, metaV1.GetOptions{})
	if err != nil {
		return controller.LogSources{}, err
	}
	return controller.LogSources{
		ContainerNames:     common.GetContainerNames(&pod.Spec),
		InitContainerNames: common.GetInitContainerNames(&pod.Spec),
		PodNames:           []string{resourceName},
	}, nil
}

func getLogSourcesFromController(kubernetes kubernetes.Interface, ns, resourceName, resourceType string) (controller.LogSources, error) {
	ref := metaV1.OwnerReference{Kind: resourceType, Name: resourceName}
	rc, err := controller.NewResourceController(ref, ns, kubernetes)
	if err != nil {
		return controller.LogSources{}, err
	}
	allPods, err := kubernetes.CoreV1().Pods(ns).List(context.TODO(), api.ListEverything)
	if err != nil {
		return controller.LogSources{}, err
	}
	return rc.GetLogSources(allPods.Items), nil
}
