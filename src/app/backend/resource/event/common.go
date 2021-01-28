package event

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
)

var EmptyEventList = &common.EventList{
	Events: make([]common.Event, 0),
	ListMeta: api.ListMeta{
		TotalItems: 0,
	},
}

func GetEvents(kubernetes kubernetes.Interface, namespace, resourceName string) ([]v1.Event, error) {
	fieldSelector, err := fields.ParseSelector("involvedObject.name" + "=" + resourceName)
	if err != nil {
		return nil, err
	}

	channels := &common.ResourceChannels{
		EventList: common.GetEventListChannelWithOptions(
			kubernetes,
			common.NewSameNamespaceQuery(namespace),
			metaV1.ListOptions{
				LabelSelector: labels.Everything().String(),
				FieldSelector: fieldSelector.String(),
			},
			1),
	}

	eventList := <-channels.EventList.List
	if err := <-channels.EventList.Error; err != nil {
		return nil, err
	}
	return FillEventsType(eventList.Items), nil
}

func GetPodsEvents(client kubernetes.Interface, namespace string, pods []v1.Pod) (
	[]v1.Event, error) {

	nsQuery := common.NewSameNamespaceQuery(namespace)
	if namespace == v1.NamespaceAll {
		nsQuery = common.NewNamespaceQuery([]string{})
	}

	channels := &common.ResourceChannels{
		EventList: common.GetEventListChannel(client, nsQuery, 1),
	}

	eventList := <-channels.EventList.List
	if err := <-channels.EventList.Error; err != nil {
		return nil, err
	}

	events := filterEventsByPodsUID(eventList.Items, pods)
	return events, nil
}

func GetPodEvents(client kubernetes.Interface, namespace, podName string) ([]v1.Event, error) {

	channels := &common.ResourceChannels{
		PodList: common.GetPodListChannel(client,
			common.NewSameNamespaceQuery(namespace),
			1),
		EventList: common.GetEventListChannel(client, common.NewSameNamespaceQuery(namespace), 1),
	}

	podList := <-channels.PodList.List
	if err := <-channels.PodList.Error; err != nil {
		return nil, err
	}

	eventList := <-channels.EventList.List
	if err := <-channels.EventList.Error; err != nil {
		return nil, err
	}

	l := make([]v1.Pod, 0)
	for _, pi := range podList.Items {
		if pi.Name == podName {
			l = append(l, pi)
		}
	}

	events := filterEventsByPodsUID(eventList.Items, l)
	return FillEventsType(events), nil
}

func GetNodeEvents(client kubernetes.Interface, dsQuery *dataselect.DataSelectQuery, nodeName string) (*common.EventList, error) {
	eventList := common.EventList{
		Events: make([]common.Event, 0),
	}

	scheme := runtime.NewScheme()
	groupVersion := schema.GroupVersion{Group: "", Version: "v1"}
	scheme.AddKnownTypes(groupVersion, &v1.Node{})

	mc := client.CoreV1().Nodes()
	node, err := mc.Get(context.TODO(), nodeName, metaV1.GetOptions{})
	if err != nil {
		return &eventList, err
	}

	events, err := client.CoreV1().Events(v1.NamespaceAll).Search(scheme, node)
	_, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return &eventList, criticalError
	}

	eventList = CreateEventList(FillEventsType(events.Items), dsQuery)
	return &eventList, nil
}

func GetNamespaceEvents(client kubernetes.Interface, dsQuery *dataselect.DataSelectQuery, namespace string) (common.EventList, error) {
	events, _ := client.CoreV1().Events(namespace).List(context.TODO(), api.ListEverything)
	return CreateEventList(FillEventsType(events.Items), dsQuery), nil
}

func FillEventsType(events []v1.Event) []v1.Event {
	for i := range events {
		if len(events[i].Type) == 0 {
			if isFailedReason(events[i].Reason, FailedReasonPartials...) {
				events[i].Type = v1.EventTypeWarning
			} else {
				events[i].Type = v1.EventTypeNormal
			}
		}
	}
	return events
}

func ToEvent(event v1.Event) common.Event {
	result := common.Event{
		ObjectMeta:      api.NewObjectMeta(event.ObjectMeta),
		TypeMeta:        api.NewTypeMeta(api.ResourceKindEvent),
		Message:         event.Message,
		SourceComponent: event.Source.Component,
		SourceHost:      event.Source.Host,
		SubObject:       event.InvolvedObject.FieldPath,
		Count:           event.Count,
		FirstSeen:       event.FirstTimestamp,
		LastSeen:        event.LastTimestamp,
		Reason:          event.Reason,
		Type:            event.Type,
	}
	return result
}

func GetResourceEvents(client kubernetes.Interface, dsQuery *dataselect.DataSelectQuery, namespace, name string) (
	*common.EventList, error) {
	resourceEvents, err := GetEvents(client, namespace, name)
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return EmptyEventList, err
	}

	events := CreateEventList(resourceEvents, dsQuery)
	events.Errors = nonCriticalErrors
	return &events, nil
}

func CreateEventList(events []v1.Event, dsQuery *dataselect.DataSelectQuery) common.EventList {
	eventList := common.EventList{
		Events:   make([]common.Event, 0),
		ListMeta: api.ListMeta{TotalItems: len(events)},
	}

	events = fromCells(dataselect.GenericDataSelect(toCells(events), dsQuery))
	for _, event := range events {
		eventDetail := ToEvent(event)
		eventList.Events = append(eventList.Events, eventDetail)
	}

	return eventList
}

type EventCell v1.Event

func (eventCell EventCell) GetProperty(name dataselect.PropertyName) dataselect.ComparableValue {
	switch name {
	case dataselect.NameProperty:
		return dataselect.StdComparableString(eventCell.ObjectMeta.Name)
	case dataselect.CreationTimestampProperty:
		return dataselect.StdComparableTime(eventCell.ObjectMeta.CreationTimestamp.Time)
	case dataselect.NamespaceProperty:
		return dataselect.StdComparableString(eventCell.ObjectMeta.Namespace)
	default:
		return nil
	}
}

func toCells(std []v1.Event) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = EventCell(std[i])
	}
	return cells
}

func fromCells(cells []dataselect.DataCell) []v1.Event {
	std := make([]v1.Event, len(cells))
	for i := range std {
		std[i] = v1.Event(cells[i].(EventCell))
	}
	return std
}
