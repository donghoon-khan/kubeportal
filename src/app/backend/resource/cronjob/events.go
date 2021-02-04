package cronjob

import (
	"k8s.io/client-go/kubernetes"

	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/event"
)

func GetCronJobEvents(kubernetes kubernetes.Interface, dsQuery *dataselect.DataSelectQuery, namespace, name string) (
	*common.EventList, error) {
	raw, err := event.GetEvents(kubernetes, namespace, name)
	if err != nil {
		return event.EmptyEventList, err
	}

	events := event.CreateEventList(raw, dsQuery)
	return &events, nil
}
