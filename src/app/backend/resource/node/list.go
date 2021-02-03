package node

import (
	"context"
	"log"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	metricApi "github.com/donghoon-khan/kubeportal/src/app/backend/integration/metric/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
)

type NodeList struct {
	ListMeta          api.ListMeta       `json:"listMeta"`
	Nodes             []Node             `json:"nodes"`
	CumulativeMetrics []metricApi.Metric `json:"cumulativeMetrics"`

	Errors []error `json:"errors"`
}

type Node struct {
	ObjectMeta         api.ObjectMeta         `json:"objectMeta"`
	TypeMeta           api.TypeMeta           `json:"typeMeta"`
	Ready              v1.ConditionStatus     `json:"ready"`
	AllocatedResources NodeAllocatedResources `json:"allocatedResources"`
}

func GetNodeList(kubernetes kubernetes.Interface, dsQuery *dataselect.DataSelectQuery, metricClient metricApi.MetricClient) (*NodeList, error) {
	nodes, err := kubernetes.CoreV1().Nodes().List(context.TODO(), api.ListEverything)

	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	return toNodeList(kubernetes, nodes.Items, nonCriticalErrors, dsQuery, metricClient), nil
}

func toNodeList(client kubernetes.Interface, nodes []v1.Node, nonCriticalErrors []error, dsQuery *dataselect.DataSelectQuery,
	metricClient metricApi.MetricClient) *NodeList {
	nodeList := &NodeList{
		Nodes:    make([]Node, 0),
		ListMeta: api.ListMeta{TotalItems: len(nodes)},
		Errors:   nonCriticalErrors,
	}

	nodeCells, metricPromises, filteredTotal := dataselect.GenericDataSelectWithFilterAndMetrics(toCells(nodes),
		dsQuery, metricApi.NoResourceCache, metricClient)
	nodes = fromCells(nodeCells)
	nodeList.ListMeta = api.ListMeta{TotalItems: filteredTotal}

	for _, node := range nodes {
		pods, err := getNodePods(client, node)
		if err != nil {
			log.Printf("Couldn't get pods of %s node: %s\n", node.Name, err)
		}

		nodeList.Nodes = append(nodeList.Nodes, toNode(node, pods))
	}

	cumulativeMetrics, err := metricPromises.GetMetrics()
	nodeList.CumulativeMetrics = cumulativeMetrics
	if err != nil {
		nodeList.CumulativeMetrics = make([]metricApi.Metric, 0)
	}

	return nodeList
}

func toNode(node v1.Node, pods *v1.PodList) Node {
	allocatedResources, err := getNodeAllocatedResources(node, pods)
	if err != nil {
		log.Printf("Couldn't get allocated resources of %s node: %s\n", node.Name, err)
	}

	return Node{
		ObjectMeta:         api.NewObjectMeta(node.ObjectMeta),
		TypeMeta:           api.NewTypeMeta(api.ResourceKindNode),
		Ready:              getNodeConditionStatus(node, v1.NodeReady),
		AllocatedResources: allocatedResources,
	}
}

func getNodeConditionStatus(node v1.Node, conditionType v1.NodeConditionType) v1.ConditionStatus {
	for _, condition := range node.Status.Conditions {
		if condition.Type == conditionType {
			return condition.Status
		}
	}
	return v1.ConditionUnknown
}
