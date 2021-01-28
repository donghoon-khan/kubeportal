package pod

import (
	"k8s.io/client-go/kubernetes"

	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/event"
)

func GetEventsForPod(kubernetes kubernetes.Interface, dsQuery *dataselect.DataSelectQuery, namespace,
	podName string) (*common.EventList, error) {
	return event.GetResourceEvents(kubernetes, dsQuery, namespace, podName)
}
