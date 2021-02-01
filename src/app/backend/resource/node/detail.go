package node

import (
	"context"
	"log"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	metricApi "github.com/donghoon-khan/kubeportal/src/app/backend/integration/metric/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/event"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/pod"
)

type NodeAllocatedResources struct {
	CPURequests            int64   `json:"cpuRequests"`
	CPURequestsFraction    float64 `json:"cpuRequestsFraction"`
	CPULimits              int64   `json:"cpuLimits"`
	CPULimitsFraction      float64 `json:"cpuLimitsFraction"`
	CPUCapacity            int64   `json:"cpuCapacity"`
	MemoryRequests         int64   `json:"memoryRequests"`
	MemoryRequestsFraction float64 `json:"memoryRequestsFraction"`
	MemoryLimits           int64   `json:"memoryLimits"`
	MemoryLimitsFraction   float64 `json:"memoryLimitsFraction"`
	MemoryCapacity         int64   `json:"memoryCapacity"`
	AllocatedPods          int     `json:"allocatedPods"`
	PodCapacity            int64   `json:"podCapacity"`
	PodFraction            float64 `json:"podFraction"`
}

type NodeDetail struct {
	Node            `json:",inline"`
	Phase           v1.NodePhase       `json:"phase"`
	PodCIDR         string             `json:"podCIDR"`
	ProviderID      string             `json:"providerID"`
	Unschedulable   bool               `json:"unschedulable"`
	NodeInfo        v1.NodeSystemInfo  `json:"nodeInfo"`
	Conditions      []common.Condition `json:"conditions"`
	ContainerImages []string           `json:"containerImages"`
	PodList         pod.PodList        `json:"podList"`
	EventList       common.EventList   `json:"eventList"`
	Metrics         []metricApi.Metric `json:"metrics"`
	Taints          []v1.Taint         `json:"taints,omitempty"`
	Addresses       []v1.NodeAddress   `json:"addresses,omitempty"`
	Errors          []error            `json:"errors" swaggertype:"array,string"`
}

func GetNodeDetail(kubernetes kubernetes.Interface, metricClient metricApi.MetricClient, name string,
	dsQuery *dataselect.DataSelectQuery) (*NodeDetail, error) {
	log.Printf("Getting details of %s node", name)

	node, err := kubernetes.CoreV1().Nodes().Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	_, metricPromises := dataselect.GenericDataSelectWithMetrics(toCells([]v1.Node{*node}),
		dsQuery,
		metricApi.NoResourceCache, metricClient)

	pods, err := getNodePods(kubernetes, *node)
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	podList, err := GetNodePods(kubernetes, metricClient, dsQuery, name)
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	eventList, err := event.GetNodeEvents(kubernetes, dsQuery, node.Name)
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	allocatedResources, err := getNodeAllocatedResources(*node, pods)
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	metrics, _ := metricPromises.GetMetrics()
	nodeDetails := toNodeDetail(*node, podList, eventList, allocatedResources, metrics, nonCriticalErrors)
	return &nodeDetails, nil
}

func getNodeAllocatedResources(node v1.Node, podList *v1.PodList) (NodeAllocatedResources, error) {
	reqs, limits := map[v1.ResourceName]resource.Quantity{}, map[v1.ResourceName]resource.Quantity{}

	for _, pod := range podList.Items {
		podReqs, podLimits, err := PodRequestsAndLimits(&pod)
		if err != nil {
			return NodeAllocatedResources{}, err
		}
		for podReqName, podReqValue := range podReqs {
			if value, ok := reqs[podReqName]; !ok {
				reqs[podReqName] = podReqValue.DeepCopy()
			} else {
				value.Add(podReqValue)
				reqs[podReqName] = value
			}
		}
		for podLimitName, podLimitValue := range podLimits {
			if value, ok := limits[podLimitName]; !ok {
				limits[podLimitName] = podLimitValue.DeepCopy()
			} else {
				value.Add(podLimitValue)
				limits[podLimitName] = value
			}
		}
	}

	cpuRequests, cpuLimits, memoryRequests, memoryLimits := reqs[v1.ResourceCPU],
		limits[v1.ResourceCPU], reqs[v1.ResourceMemory], limits[v1.ResourceMemory]

	var cpuRequestsFraction, cpuLimitsFraction float64 = 0, 0
	if capacity := float64(node.Status.Capacity.Cpu().MilliValue()); capacity > 0 {
		cpuRequestsFraction = float64(cpuRequests.MilliValue()) / capacity * 100
		cpuLimitsFraction = float64(cpuLimits.MilliValue()) / capacity * 100
	}

	var memoryRequestsFraction, memoryLimitsFraction float64 = 0, 0
	if capacity := float64(node.Status.Capacity.Memory().MilliValue()); capacity > 0 {
		memoryRequestsFraction = float64(memoryRequests.MilliValue()) / capacity * 100
		memoryLimitsFraction = float64(memoryLimits.MilliValue()) / capacity * 100
	}

	var podFraction float64 = 0
	var podCapacity int64 = node.Status.Capacity.Pods().Value()
	if podCapacity > 0 {
		podFraction = float64(len(podList.Items)) / float64(podCapacity) * 100
	}

	return NodeAllocatedResources{
		CPURequests:            cpuRequests.MilliValue(),
		CPURequestsFraction:    cpuRequestsFraction,
		CPULimits:              cpuLimits.MilliValue(),
		CPULimitsFraction:      cpuLimitsFraction,
		CPUCapacity:            node.Status.Capacity.Cpu().MilliValue(),
		MemoryRequests:         memoryRequests.Value(),
		MemoryRequestsFraction: memoryRequestsFraction,
		MemoryLimits:           memoryLimits.Value(),
		MemoryLimitsFraction:   memoryLimitsFraction,
		MemoryCapacity:         node.Status.Capacity.Memory().Value(),
		AllocatedPods:          len(podList.Items),
		PodCapacity:            podCapacity,
		PodFraction:            podFraction,
	}, nil
}

func PodRequestsAndLimits(pod *v1.Pod) (reqs, limits v1.ResourceList, err error) {
	reqs, limits = v1.ResourceList{}, v1.ResourceList{}
	for _, container := range pod.Spec.Containers {
		addResourceList(reqs, container.Resources.Requests)
		addResourceList(limits, container.Resources.Limits)
	}

	for _, container := range pod.Spec.InitContainers {
		maxResourceList(reqs, container.Resources.Requests)
		maxResourceList(limits, container.Resources.Limits)
	}

	if pod.Spec.Overhead != nil {
		addResourceList(reqs, pod.Spec.Overhead)

		for name, quantity := range pod.Spec.Overhead {
			if value, ok := limits[name]; ok && !value.IsZero() {
				value.Add(quantity)
				limits[name] = value
			}
		}
	}
	return
}

func addResourceList(list, new v1.ResourceList) {
	for name, quantity := range new {
		if value, ok := list[name]; !ok {
			list[name] = quantity.DeepCopy()
		} else {
			value.Add(quantity)
			list[name] = value
		}
	}
}

func maxResourceList(list, new v1.ResourceList) {
	for name, quantity := range new {
		if value, ok := list[name]; !ok {
			list[name] = quantity.DeepCopy()
			continue
		} else {
			if quantity.Cmp(value) > 0 {
				list[name] = quantity.DeepCopy()
			}
		}
	}
}

func GetNodePods(kubernetes kubernetes.Interface, metricClient metricApi.MetricClient,
	dsQuery *dataselect.DataSelectQuery, name string) (*pod.PodList, error) {
	podList := pod.PodList{
		Pods:              []pod.Pod{},
		CumulativeMetrics: []metricApi.Metric{},
	}

	node, err := kubernetes.CoreV1().Nodes().Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return &podList, err
	}

	pods, err := getNodePods(kubernetes, *node)
	podNonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return &podList, criticalError
	}

	events, err := event.GetPodsEvents(kubernetes, v1.NamespaceAll, pods.Items)
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return &podList, criticalError
	}

	nonCriticalErrors = append(nonCriticalErrors, podNonCriticalErrors...)
	podList = pod.ToPodList(pods.Items, events, nonCriticalErrors, dsQuery, metricClient)
	return &podList, nil
}

func getNodePods(kubernetes kubernetes.Interface, node v1.Node) (*v1.PodList, error) {
	fieldSelector, err := fields.ParseSelector("spec.nodeName=" + node.Name +
		",status.phase!=" + string(v1.PodSucceeded) +
		",status.phase!=" + string(v1.PodFailed))

	if err != nil {
		return nil, err
	}

	return kubernetes.CoreV1().Pods(v1.NamespaceAll).List(context.TODO(), metaV1.ListOptions{
		FieldSelector: fieldSelector.String(),
	})
}

func toNodeDetail(node v1.Node, pods *pod.PodList, eventList *common.EventList,
	allocatedResources NodeAllocatedResources, metrics []metricApi.Metric, nonCriticalErrors []error) NodeDetail {
	return NodeDetail{
		Node: Node{
			ObjectMeta:         api.NewObjectMeta(node.ObjectMeta),
			TypeMeta:           api.NewTypeMeta(api.ResourceKindNode),
			AllocatedResources: allocatedResources,
		},
		Phase:           node.Status.Phase,
		ProviderID:      node.Spec.ProviderID,
		PodCIDR:         node.Spec.PodCIDR,
		Unschedulable:   node.Spec.Unschedulable,
		NodeInfo:        node.Status.NodeInfo,
		Conditions:      getNodeConditions(node),
		ContainerImages: getContainerImages(node),
		PodList:         *pods,
		EventList:       *eventList,
		Metrics:         metrics,
		Taints:          node.Spec.Taints,
		Addresses:       node.Status.Addresses,
		Errors:          nonCriticalErrors,
	}
}
