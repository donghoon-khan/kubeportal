package common

import api "k8s.io/api/core/v1"

type PodInfo struct {
	Current   int32   `json:"current"`
	Desired   *int32  `json:"desired,omitempty"`
	Running   int32   `json:"running"`
	Pending   int32   `json:"pending"`
	Failed    int32   `json:"failed"`
	Succeeded int32   `json:"succeeded"`
	Warnings  []Event `json:"warnings"`
}

func GetPodInfo(current int32, desired *int32, pods []api.Pod) PodInfo {
	result := PodInfo{
		Current:  current,
		Desired:  desired,
		Warnings: make([]Event, 0),
	}

	for _, pod := range pods {
		switch pod.Status.Phase {
		case api.PodRunning:
			result.Running++
		case api.PodPending:
			result.Pending++
		case api.PodFailed:
			result.Failed++
		case api.PodSucceeded:
			result.Succeeded++
		}
	}

	return result
}
