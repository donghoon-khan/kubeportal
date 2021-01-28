package pod

/*
import (
	"log"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	metricApi "github.com/donghoon-khan/kubeportal/src/app/backend/integration/metric/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/event"
)

type PodList struct {
	ListMeta          api.ListMeta          `json:"listMeta"`
	CumulativeMetrics []metricApi.Metric    `json:"cumulativeMetrics"`
	Status            common.ResourceStatus `json:"status"`
	Pods              []Pod                 `json:"pods"`
	Errors            []error               `json:"errors"`
}

type PodStatus struct {
	Status          string              `json:"status"`
	PodPhase        v1.PodPhase         `json:"podPhase"`
	ContainerStates []v1.ContainerState `json:"containerStates"`
}

type Pod struct {
	ObjectMeta   api.ObjectMeta `json:"objectMeta"`
	TypeMeta     api.TypeMeta   `json:"typeMeta"`
	Status       string         `json:"status"`
	RestartCount int32          `json:"restartCount"`
	Metrics      *PodMetrics    `json:"metrics"`
	Warnings     []common.Event `json:"warnings"`
	NodeName     string         `json:"nodeName"`
}

var EmptyPodList = &PodList{
	Pods:   make([]Pod, 0),
	Errors: make([]error, 0),
	ListMeta: api.ListMeta{
		TotalItems: 0,
	},
}

func GetPodList(kubernetes kubernetes.Interface, metricClient metricApi.MetricClient, nsQuery *common.NamespaceQuery,
	dsQuery *dataselect.DataSelectQuery) (*PodList, error) {
	log.Print("Getting list of all pods in the cluster")

	channels := &common.ResourceChannels{
		PodList:   common.GetPodListChannelWithOptions(kubernetes, nsQuery, metaV1.ListOptions{}, 1),
		EventList: common.GetEventListChannel(kubernetes, nsQuery, 1),
	}

	return GetPodListFromChannels(channels, dsQuery, metricClient)
}

func GetPodListFromChannels(channels *common.ResourceChannels, dsQuery *dataselect.DataSelectQuery,
	metricClient metricApi.MetricClient) (*PodList, error) {

	pods := <-channels.PodList.List
	err := <-channels.PodList.Error
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	eventList := <-channels.EventList.List
	err = <-channels.EventList.Error
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	podList := ToPodList(pods.Items, eventList.Items, nonCriticalErrors, dsQuery, metricClient)
	podList.Status = getStatus(pods, eventList.Items)
	return &podList, nil
}

func ToPodList(pods []v1.Pod, events []v1.Event, nonCriticalErrors []error, dsQuery *dataselect.DataSelectQuery,
	metricClient metricApi.MetricClient) PodList {
	podList := PodList{
		Pods:   make([]Pod, 0),
		Errors: nonCriticalErrors,
	}

	podCells, cumulativeMetricsPromises, filteredTotal := dataselect.
		GenericDataSelectWithFilterAndMetrics(toCells(pods), dsQuery, metricApi.NoResourceCache, metricClient)
	pods = fromCells(podCells)
	podList.ListMeta = api.ListMeta{TotalItems: filteredTotal}

	metrics, err := getMetricsPerPod(pods, metricClient, dsQuery)
	if err != nil {
		log.Printf("Skipping metrics because of error: %s\n", err)
	}

	for _, pod := range pods {
		warnings := event.GetPodsEventWarnings(events, []v1.Pod{pod})
		podDetail := toPod(&pod, metrics, warnings)
		podList.Pods = append(podList.Pods, podDetail)
	}

	cumulativeMetrics, err := cumulativeMetricsPromises.GetMetrics()
	if err != nil {
		log.Printf("Skipping metrics because of error: %s\n", err)
		cumulativeMetrics = make([]metricApi.Metric, 0)
	}

	podList.CumulativeMetrics = cumulativeMetrics
	return podList
}

func toPod(pod *v1.Pod, metrics *MetricsByPod, warnings []common.Event) Pod {
	podDetail := Pod{
		ObjectMeta:   api.NewObjectMeta(pod.ObjectMeta),
		TypeMeta:     api.NewTypeMeta(api.ResourceKindPod),
		Warnings:     warnings,
		Status:       getPodStatus(*pod),
		RestartCount: getRestartCount(*pod),
		NodeName:     pod.Spec.NodeName,
	}

	if m, exists := metrics.MetricsMap[pod.UID]; exists {
		podDetail.Metrics = &m
	}

	return podDetail
}
*/
