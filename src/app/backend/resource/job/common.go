package job

import (
	batch "k8s.io/api/batch/v1"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	metricApi "github.com/donghoon-khan/kubeportal/src/app/backend/integration/metric/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
)

type JobCell batch.Job

func (jobCell JobCell) GetProperty(name dataselect.PropertyName) dataselect.ComparableValue {
	switch name {
	case dataselect.NameProperty:
		return dataselect.StdComparableString(jobCell.ObjectMeta.Name)
	case dataselect.CreationTimestampProperty:
		return dataselect.StdComparableTime(jobCell.ObjectMeta.CreationTimestamp.Time)
	case dataselect.NamespaceProperty:
		return dataselect.StdComparableString(jobCell.ObjectMeta.Namespace)
	default:
		return nil
	}
}

func (jobCell JobCell) GetResourceSelector() *metricApi.ResourceSelector {
	return &metricApi.ResourceSelector{
		Namespace:    jobCell.ObjectMeta.Namespace,
		ResourceType: api.ResourceKindJob,
		ResourceName: jobCell.ObjectMeta.Name,
		UID:          jobCell.UID,
	}
}

func ToCell(std []batch.Job) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = JobCell(std[i])
	}
	return cells
}

func FromCells(cells []dataselect.DataCell) []batch.Job {
	std := make([]batch.Job, len(cells))
	for i := range std {
		std[i] = batch.Job(cells[i].(JobCell))
	}
	return std
}

/*func getStatus(list *batch.JobList, pods []v1.Pod) common.ResourceStatus {
	info := common.ResourceStatus{}
	if list == nil {
		return info
	}

	for _, job := range list.Items {
		matchingPods := common.FilterPodsForJob(job, pods)
		podInfo := common.GetPodInfo(job.Status.Active, job.Spec.Completions, matchingPods)
		jobStatus := getJobStatus(&job)

		if jobStatus.Status == JobStatusFailed {
			info.Failed++
		} else if jobStatus.Status == JobStatusComplete {
			info.Succeeded++
		} else if podInfo.Running > 0 {
			info.Running++
		} else {
			info.Pending++
		}
	}

	return info
}
*/
