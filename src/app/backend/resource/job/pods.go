package job

import (
	"context"
	"log"

	batch "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
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

func GetJobPods(kubernetes kubernetes.Interface, metricClient metricApi.MetricClient,
	dsQuery *dataselect.DataSelectQuery, namespace string, jobName string) (*pod.PodList, error) {
	log.Printf("Getting replication controller %s pods in namespace %s", jobName, namespace)

	pods, err := getRawJobPods(kubernetes, jobName, namespace)
	if err != nil {
		return pod.EmptyPodList, err
	}

	events, err := event.GetPodsEvents(kubernetes, namespace, pods)
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return pod.EmptyPodList, criticalError
	}

	podList := pod.ToPodList(pods, events, nonCriticalErrors, dsQuery, metricClient)
	return &podList, nil
}

func getRawJobPods(kubernetes kubernetes.Interface, petSetName, namespace string) ([]v1.Pod, error) {
	job, err := kubernetes.BatchV1().Jobs(namespace).Get(context.TODO(), petSetName, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	labelSelector := labels.SelectorFromSet(job.Spec.Selector.MatchLabels)
	channels := &common.ResourceChannels{
		PodList: common.GetPodListChannelWithOptions(kubernetes, common.NewSameNamespaceQuery(namespace),
			metaV1.ListOptions{
				LabelSelector: labelSelector.String(),
				FieldSelector: fields.Everything().String(),
			}, 1),
	}

	podList := <-channels.PodList.List
	if err := <-channels.PodList.Error; err != nil {
		return nil, err
	}

	return podList.Items, nil
}

func getJobPodInfo(kubernetes kubernetes.Interface, job *batch.Job) (*common.PodInfo, error) {
	labelSelector := labels.SelectorFromSet(job.Spec.Selector.MatchLabels)
	channels := &common.ResourceChannels{
		PodList: common.GetPodListChannelWithOptions(kubernetes, common.NewSameNamespaceQuery(
			job.Namespace),
			metaV1.ListOptions{
				LabelSelector: labelSelector.String(),
				FieldSelector: fields.Everything().String(),
			}, 1),
	}

	pods := <-channels.PodList.List
	if err := <-channels.PodList.Error; err != nil {
		return nil, err
	}

	podInfo := common.GetPodInfo(job.Status.Active, job.Spec.Completions, pods.Items)
	podInfo.Running = job.Status.Active
	podInfo.Succeeded = job.Status.Succeeded
	podInfo.Failed = job.Status.Failed

	return &podInfo, nil
}
