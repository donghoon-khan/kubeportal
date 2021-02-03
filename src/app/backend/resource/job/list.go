package job

import (
	"log"

	batch "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	metricApi "github.com/donghoon-khan/kubeportal/src/app/backend/integration/metric/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/event"
)

type JobList struct {
	ListMeta          api.ListMeta          `json:"listMeta"`
	CumulativeMetrics []metricApi.Metric    `json:"cumulativeMetrics"`
	Status            common.ResourceStatus `json:"status"`
	Jobs              []Job                 `json:"jobs"`
	Errors            []error               `json:"errors"`
}

type JobStatusType string

const (
	JobStatusRunning  JobStatusType = "Running"
	JobStatusComplete JobStatusType = "Complete"
	JobStatusFailed   JobStatusType = "Failed"
)

type JobStatus struct {
	Status     JobStatusType      `json:"status"`
	Message    string             `json:"message"`
	Conditions []common.Condition `json:"conditions"`
}

type Job struct {
	ObjectMeta api.ObjectMeta `json:"objectMeta"`
	TypeMeta   api.TypeMeta   `json:"typeMeta"`

	Pods                common.PodInfo `json:"podInfo"`
	ContainerImages     []string       `json:"containerImages"`
	InitContainerImages []string       `json:"initContainerImages"`
	Parallelism         *int32         `json:"parallelism"`
	JobStatus           JobStatus      `json:"jobStatus"`
}

func GetJobList(kubernetes kubernetes.Interface, nsQuery *common.NamespaceQuery,
	dsQuery *dataselect.DataSelectQuery, metricClient metricApi.MetricClient) (*JobList, error) {
	log.Print("Getting list of all jobs in the cluster")

	channels := &common.ResourceChannels{
		JobList:   common.GetJobListChannel(kubernetes, nsQuery, 1),
		PodList:   common.GetPodListChannel(kubernetes, nsQuery, 1),
		EventList: common.GetEventListChannel(kubernetes, nsQuery, 1),
	}

	return GetJobListFromChannels(channels, dsQuery, metricClient)
}

func GetJobListFromChannels(channels *common.ResourceChannels, dsQuery *dataselect.DataSelectQuery,
	metricClient metricApi.MetricClient) (*JobList, error) {

	jobs := <-channels.JobList.List
	err := <-channels.JobList.Error
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	pods := <-channels.PodList.List
	err = <-channels.PodList.Error
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	events := <-channels.EventList.List
	err = <-channels.EventList.Error
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	jobList := ToJobList(jobs.Items, pods.Items, events.Items, nonCriticalErrors, dsQuery, metricClient)
	jobList.Status = getStatus(jobs, pods.Items)
	return jobList, nil
}

func ToJobList(jobs []batch.Job, pods []v1.Pod, events []v1.Event, nonCriticalErrors []error,
	dsQuery *dataselect.DataSelectQuery, metricClient metricApi.MetricClient) *JobList {

	jobList := &JobList{
		Jobs:     make([]Job, 0),
		ListMeta: api.ListMeta{TotalItems: len(jobs)},
		Errors:   nonCriticalErrors,
	}

	cachedResources := &metricApi.CachedResources{
		Pods: pods,
	}
	jobCells, metricPromises, filteredTotal := dataselect.GenericDataSelectWithFilterAndMetrics(ToCells(jobs),
		dsQuery, cachedResources, metricClient)
	jobs = FromCells(jobCells)
	jobList.ListMeta = api.ListMeta{TotalItems: filteredTotal}

	for _, job := range jobs {
		matchingPods := common.FilterPodsForJob(job, pods)
		podInfo := common.GetPodInfo(job.Status.Active, job.Spec.Completions, matchingPods)
		podInfo.Warnings = event.GetPodsEventWarnings(events, matchingPods)
		jobList.Jobs = append(jobList.Jobs, toJob(&job, &podInfo))
	}

	cumulativeMetrics, err := metricPromises.GetMetrics()
	jobList.CumulativeMetrics = cumulativeMetrics
	if err != nil {
		jobList.CumulativeMetrics = make([]metricApi.Metric, 0)
	}

	return jobList
}

func toJob(job *batch.Job, podInfo *common.PodInfo) Job {
	return Job{
		ObjectMeta:          api.NewObjectMeta(job.ObjectMeta),
		TypeMeta:            api.NewTypeMeta(api.ResourceKindJob),
		ContainerImages:     common.GetContainerImages(&job.Spec.Template.Spec),
		InitContainerImages: common.GetInitContainerImages(&job.Spec.Template.Spec),
		Pods:                *podInfo,
		JobStatus:           getJobStatus(job),
		Parallelism:         job.Spec.Parallelism,
	}
}

func getJobStatus(job *batch.Job) JobStatus {
	jobStatus := JobStatus{Status: JobStatusRunning, Conditions: getJobConditions(job)}
	for _, condition := range job.Status.Conditions {
		if condition.Type == batch.JobComplete && condition.Status == v1.ConditionTrue {
			jobStatus.Status = JobStatusComplete
			break
		} else if condition.Type == batch.JobFailed && condition.Status == v1.ConditionTrue {
			jobStatus.Status = JobStatusFailed
			jobStatus.Message = condition.Message
			break
		}
	}
	return jobStatus
}

func getJobConditions(job *batch.Job) []common.Condition {
	var conditions []common.Condition
	for _, condition := range job.Status.Conditions {
		conditions = append(conditions, common.Condition{
			Type:               string(condition.Type),
			Status:             condition.Status,
			LastProbeTime:      condition.LastProbeTime,
			LastTransitionTime: condition.LastTransitionTime,
			Reason:             condition.Reason,
			Message:            condition.Message,
		})
	}
	return conditions
}
