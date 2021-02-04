package service

import (
	"context"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"

	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	metricApi "github.com/donghoon-khan/kubeportal/src/app/backend/integration/metric/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/event"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/pod"
)

func GetServicePods(kubernetes kubernetes.Interface, metricClient metricApi.MetricClient, namespace,
	name string, dsQuery *dataselect.DataSelectQuery) (*pod.PodList, error) {
	podList := pod.PodList{
		Pods:              []pod.Pod{},
		CumulativeMetrics: []metricApi.Metric{},
	}

	service, err := kubernetes.CoreV1().Services(namespace).Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return &podList, err
	}

	if service.Spec.Selector == nil {
		return &podList, nil
	}

	labelSelector := labels.SelectorFromSet(service.Spec.Selector)
	channels := &common.ResourceChannels{
		PodList: common.GetPodListChannelWithOptions(kubernetes, common.NewSameNamespaceQuery(namespace),
			metaV1.ListOptions{
				LabelSelector: labelSelector.String(),
				FieldSelector: fields.Everything().String(),
			}, 1),
	}

	apiPodList := <-channels.PodList.List
	if err := <-channels.PodList.Error; err != nil {
		return &podList, err
	}

	events, err := event.GetPodsEvents(kubernetes, namespace, apiPodList.Items)
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return &podList, criticalError
	}

	podList = pod.ToPodList(apiPodList.Items, events, nonCriticalErrors, dsQuery, metricClient)
	return &podList, nil
}
