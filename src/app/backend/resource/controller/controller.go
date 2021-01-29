package controller

import (
	"context"
	"fmt"
	"strings"

	apps "k8s.io/api/apps/v1"
	batch "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/event"
)

type ResourceOwner struct {
	ObjectMeta          api.ObjectMeta `json:"objectMeta"`
	TypeMeta            api.TypeMeta   `json:"typeMeta"`
	Pods                common.PodInfo `json:"pods"`
	ContainerImages     []string       `json:"containerImages"`
	InitContainerImages []string       `json:"initContainerImages"`
}

type LogSources struct {
	ContainerNames     []string `json:"containerNames"`
	InitContainerNames []string `json:"initContainerNames"`
	PodNames           []string `json:"podNames"`
}

type ResourceController interface {
	UID() types.UID
	Get(allPods []v1.Pod, allEvents []v1.Event) ResourceOwner
	GetLogSources(allPods []v1.Pod) LogSources
}

func NewResourceController(ref metaV1.OwnerReference, namespace string, kubernetes kubernetes.Interface) (
	ResourceController, error) {
	switch strings.ToLower(ref.Kind) {
	case api.ResourceKindJob:
		job, err := kubernetes.BatchV1().Jobs(namespace).Get(context.TODO(), ref.Name, metaV1.GetOptions{})
		if err != nil {
			return nil, err
		}
		return JobController(*job), nil
	case api.ResourceKindPod:
		pod, err := kubernetes.CoreV1().Pods(namespace).Get(context.TODO(), ref.Name, metaV1.GetOptions{})
		if err != nil {
			return nil, err
		}
		return PodController(*pod), nil
	case api.ResourceKindReplicaSet:
		rs, err := kubernetes.AppsV1().ReplicaSets(namespace).Get(context.TODO(), ref.Name, metaV1.GetOptions{})
		if err != nil {
			return nil, err
		}
		return ReplicaSetController(*rs), nil
	case api.ResourceKindReplicationController:
		rc, err := kubernetes.CoreV1().ReplicationControllers(namespace).Get(context.TODO(), ref.Name, metaV1.GetOptions{})
		if err != nil {
			return nil, err
		}
		return ReplicationControllerController(*rc), nil
	case api.ResourceKindDaemonSet:
		ds, err := kubernetes.AppsV1().DaemonSets(namespace).Get(context.TODO(), ref.Name, metaV1.GetOptions{})
		if err != nil {
			return nil, err
		}
		return DaemonSetController(*ds), nil
	case api.ResourceKindStatefulSet:
		ss, err := kubernetes.AppsV1().StatefulSets(namespace).Get(context.TODO(), ref.Name, metaV1.GetOptions{})
		if err != nil {
			return nil, err
		}
		return StatefulSetController(*ss), nil
	default:
		return nil, fmt.Errorf("unknown reference kind: %s", ref.Kind)
	}
}

type JobController batch.Job

func (self JobController) Get(allPods []v1.Pod, allEvents []v1.Event) ResourceOwner {
	matchingPods := common.FilterPodsForJob(batch.Job(self), allPods)
	podInfo := common.GetPodInfo(self.Status.Active, self.Spec.Completions, matchingPods)
	podInfo.Warnings = event.GetPodsEventWarnings(allEvents, matchingPods)

	return ResourceOwner{
		TypeMeta:            api.NewTypeMeta(api.ResourceKindJob),
		ObjectMeta:          api.NewObjectMeta(self.ObjectMeta),
		Pods:                podInfo,
		ContainerImages:     common.GetContainerImages(&self.Spec.Template.Spec),
		InitContainerImages: common.GetInitContainerImages(&self.Spec.Template.Spec),
	}
}

func (self JobController) UID() types.UID {
	return batch.Job(self).UID
}

func (self JobController) GetLogSources(allPods []v1.Pod) LogSources {
	controlledPods := common.FilterPodsForJob(batch.Job(self), allPods)
	return LogSources{
		PodNames:           getPodNames(controlledPods),
		ContainerNames:     common.GetContainerNames(&self.Spec.Template.Spec),
		InitContainerNames: common.GetInitContainerNames(&self.Spec.Template.Spec),
	}
}

type PodController v1.Pod

func (self PodController) Get(allPods []v1.Pod, allEvents []v1.Event) ResourceOwner {
	matchingPods := common.FilterPodsByControllerRef(&self, allPods)
	podInfo := common.GetPodInfo(int32(len(matchingPods)), nil, matchingPods)
	podInfo.Warnings = event.GetPodsEventWarnings(allEvents, matchingPods)

	return ResourceOwner{
		TypeMeta:            api.NewTypeMeta(api.ResourceKindPod),
		ObjectMeta:          api.NewObjectMeta(self.ObjectMeta),
		Pods:                podInfo,
		ContainerImages:     common.GetNonduplicateContainerImages(matchingPods),
		InitContainerImages: common.GetNonduplicateInitContainerImages(matchingPods),
	}
}

func (self PodController) UID() types.UID {
	return v1.Pod(self).UID
}

func (self PodController) GetLogSources(allPods []v1.Pod) LogSources {
	controlledPods := common.FilterPodsByControllerRef(&self, allPods)
	return LogSources{
		PodNames:           getPodNames(controlledPods),
		ContainerNames:     common.GetNonduplicateContainerNames(controlledPods),
		InitContainerNames: common.GetNonduplicateInitContainerNames(controlledPods),
	}
}

type ReplicaSetController apps.ReplicaSet

func (self ReplicaSetController) Get(allPods []v1.Pod, allEvents []v1.Event) ResourceOwner {
	matchingPods := common.FilterPodsByControllerRef(&self, allPods)
	podInfo := common.GetPodInfo(self.Status.Replicas, self.Spec.Replicas, matchingPods)
	podInfo.Warnings = event.GetPodsEventWarnings(allEvents, matchingPods)

	return ResourceOwner{
		TypeMeta:            api.NewTypeMeta(api.ResourceKindReplicaSet),
		ObjectMeta:          api.NewObjectMeta(self.ObjectMeta),
		Pods:                podInfo,
		ContainerImages:     common.GetContainerImages(&self.Spec.Template.Spec),
		InitContainerImages: common.GetInitContainerImages(&self.Spec.Template.Spec),
	}
}

func (self ReplicaSetController) UID() types.UID {
	return apps.ReplicaSet(self).UID
}

func (self ReplicaSetController) GetLogSources(allPods []v1.Pod) LogSources {
	controlledPods := common.FilterPodsByControllerRef(&self, allPods)
	return LogSources{
		PodNames:           getPodNames(controlledPods),
		ContainerNames:     common.GetContainerNames(&self.Spec.Template.Spec),
		InitContainerNames: common.GetInitContainerNames(&self.Spec.Template.Spec),
	}
}

type ReplicationControllerController v1.ReplicationController

func (self ReplicationControllerController) Get(allPods []v1.Pod,
	allEvents []v1.Event) ResourceOwner {
	matchingPods := common.FilterPodsByControllerRef(&self, allPods)
	podInfo := common.GetPodInfo(self.Status.Replicas, self.Spec.Replicas, matchingPods)
	podInfo.Warnings = event.GetPodsEventWarnings(allEvents, matchingPods)

	return ResourceOwner{
		TypeMeta:            api.NewTypeMeta(api.ResourceKindReplicationController),
		ObjectMeta:          api.NewObjectMeta(self.ObjectMeta),
		Pods:                podInfo,
		ContainerImages:     common.GetContainerImages(&self.Spec.Template.Spec),
		InitContainerImages: common.GetInitContainerImages(&self.Spec.Template.Spec),
	}
}

func (self ReplicationControllerController) UID() types.UID {
	return v1.ReplicationController(self).UID
}

func (self ReplicationControllerController) GetLogSources(allPods []v1.Pod) LogSources {
	controlledPods := common.FilterPodsByControllerRef(&self, allPods)
	return LogSources{
		PodNames:           getPodNames(controlledPods),
		ContainerNames:     common.GetContainerNames(&self.Spec.Template.Spec),
		InitContainerNames: common.GetInitContainerNames(&self.Spec.Template.Spec),
	}
}

type DaemonSetController apps.DaemonSet

func (self DaemonSetController) Get(allPods []v1.Pod, allEvents []v1.Event) ResourceOwner {
	matchingPods := common.FilterPodsByControllerRef(&self, allPods)
	podInfo := common.GetPodInfo(self.Status.CurrentNumberScheduled,
		&self.Status.DesiredNumberScheduled, matchingPods)
	podInfo.Warnings = event.GetPodsEventWarnings(allEvents, matchingPods)

	return ResourceOwner{
		TypeMeta:            api.NewTypeMeta(api.ResourceKindDaemonSet),
		ObjectMeta:          api.NewObjectMeta(self.ObjectMeta),
		Pods:                podInfo,
		ContainerImages:     common.GetContainerImages(&self.Spec.Template.Spec),
		InitContainerImages: common.GetInitContainerImages(&self.Spec.Template.Spec),
	}
}

func (self DaemonSetController) UID() types.UID {
	return apps.DaemonSet(self).UID
}

func (self DaemonSetController) GetLogSources(allPods []v1.Pod) LogSources {
	controlledPods := common.FilterPodsByControllerRef(&self, allPods)
	return LogSources{
		PodNames:           getPodNames(controlledPods),
		ContainerNames:     common.GetContainerNames(&self.Spec.Template.Spec),
		InitContainerNames: common.GetInitContainerNames(&self.Spec.Template.Spec),
	}
}

type StatefulSetController apps.StatefulSet

func (self StatefulSetController) Get(allPods []v1.Pod, allEvents []v1.Event) ResourceOwner {
	matchingPods := common.FilterPodsByControllerRef(&self, allPods)
	podInfo := common.GetPodInfo(self.Status.Replicas, self.Spec.Replicas, matchingPods)
	podInfo.Warnings = event.GetPodsEventWarnings(allEvents, matchingPods)

	return ResourceOwner{
		TypeMeta:            api.NewTypeMeta(api.ResourceKindStatefulSet),
		ObjectMeta:          api.NewObjectMeta(self.ObjectMeta),
		Pods:                podInfo,
		ContainerImages:     common.GetContainerImages(&self.Spec.Template.Spec),
		InitContainerImages: common.GetInitContainerImages(&self.Spec.Template.Spec),
	}
}

func (self StatefulSetController) UID() types.UID {
	return apps.StatefulSet(self).UID
}

func (self StatefulSetController) GetLogSources(allPods []v1.Pod) LogSources {
	controlledPods := common.FilterPodsByControllerRef(&self, allPods)
	return LogSources{
		PodNames:           getPodNames(controlledPods),
		ContainerNames:     common.GetContainerNames(&self.Spec.Template.Spec),
		InitContainerNames: common.GetInitContainerNames(&self.Spec.Template.Spec),
	}
}

func getPodNames(pods []v1.Pod) []string {
	names := make([]string, 0)
	for _, pod := range pods {
		names = append(names, pod.Name)
	}
	return names
}
