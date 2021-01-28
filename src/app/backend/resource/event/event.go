package event

import (
	"strings"

	api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
)

var FailedReasonPartials = []string{"failed", "err", "exceeded", "invalid", "unhealthy",
	"mismatch", "insufficient", "conflict", "outof", "nil", "backoff"}

func GetPodsEventWarnings(events []api.Event, pods []api.Pod) []common.Event {
	result := make([]common.Event, 0)
	events = getWarningEvents(events)
	failedPods := make([]api.Pod, 0)
	for _, pod := range pods {
		if !isReadyOrSucceeded(pod) {
			failedPods = append(failedPods, pod)
		}
	}

	events = filterEventsByPodsUID(events, failedPods)
	events = removeDuplicates(events)
	for _, event := range events {
		result = append(result, common.Event{
			Message: event.Message,
			Reason:  event.Reason,
			Type:    event.Type,
		})
	}

	return result
}

func filterEventsByPodsUID(events []api.Event, pods []api.Pod) []api.Event {
	result := make([]api.Event, 0)
	podEventMap := make(map[types.UID]bool, 0)

	if len(pods) == 0 || len(events) == 0 {
		return result
	}

	for _, pod := range pods {
		podEventMap[pod.UID] = true
	}

	for _, event := range events {
		if _, exists := podEventMap[event.InvolvedObject.UID]; exists {
			result = append(result, event)
		}
	}

	return result
}

func getWarningEvents(events []api.Event) []api.Event {
	return filterEventsByType(FillEventsType(events), api.EventTypeWarning)
}

func filterEventsByType(events []api.Event, eventType string) []api.Event {
	if len(eventType) == 0 || len(events) == 0 {
		return events
	}

	result := make([]api.Event, 0)
	for _, event := range events {
		if event.Type == eventType {
			result = append(result, event)
		}
	}

	return result
}

func isFailedReason(reason string, partials ...string) bool {
	for _, partial := range partials {
		if strings.Contains(strings.ToLower(reason), partial) {
			return true
		}
	}

	return false
}

func removeDuplicates(slice []api.Event) []api.Event {
	visited := make(map[string]bool, 0)
	result := make([]api.Event, 0)

	for _, elem := range slice {
		if !visited[elem.Reason] {
			visited[elem.Reason] = true
			result = append(result, elem)
		}
	}

	return result
}

func isReadyOrSucceeded(pod api.Pod) bool {
	if pod.Status.Phase == api.PodSucceeded {
		return true
	}

	if pod.Status.Phase == api.PodRunning {
		for _, c := range pod.Status.Conditions {
			if c.Type == api.PodReady {
				if c.Status == api.ConditionFalse {
					return false
				}
			}
		}
		return true
	}

	return false
}
