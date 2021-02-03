package cronjob

import (
	"context"

	batch "k8s.io/api/batch/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/kubernetes"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	metricApi "github.com/donghoon-khan/kubeportal/src/app/backend/integration/metric/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/job"
)

var emptyJobList = &job.JobList{
	Jobs:   make([]job.Job, 0),
	Errors: make([]error, 0),
	ListMeta: api.ListMeta{
		TotalItems: 0,
	},
}

func GetCronJobJobs(client kubernetes.Interface, metricClient metricApi.MetricClient,
	dsQuery *dataselect.DataSelectQuery, namespace, name string, active bool) (*job.JobList, error) {

	cronJob, err := client.BatchV1beta1().CronJobs(namespace).Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return emptyJobList, err
	}

	channels := &common.ResourceChannels{
		JobList:   common.GetJobListChannel(client, common.NewSameNamespaceQuery(namespace), 1),
		PodList:   common.GetPodListChannel(client, common.NewSameNamespaceQuery(namespace), 1),
		EventList: common.GetEventListChannel(client, common.NewSameNamespaceQuery(namespace), 1),
	}

	jobs := <-channels.JobList.List
	err = <-channels.JobList.Error
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return emptyJobList, nil
	}

	pods := <-channels.PodList.List
	err = <-channels.PodList.Error
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return emptyJobList, criticalError
	}

	events := <-channels.EventList.List
	err = <-channels.EventList.Error
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return emptyJobList, criticalError
	}

	jobs.Items = filterJobsByOwnerUID(cronJob.UID, jobs.Items)
	jobs.Items = filterJobsByState(active, jobs.Items)

	return job.ToJobList(jobs.Items, pods.Items, events.Items, nonCriticalErrors, dsQuery, metricClient), nil
}

func TriggerCronJob(client kubernetes.Interface,
	namespace, name string) error {

	cronJob, err := client.BatchV1beta1().CronJobs(namespace).Get(context.TODO(), name, metaV1.GetOptions{})

	if err != nil {
		return err
	}

	annotations := make(map[string]string)
	annotations["cronjob.kubernetes.io/instantiate"] = "manual"

	labels := make(map[string]string)
	for k, v := range cronJob.Spec.JobTemplate.Labels {
		labels[k] = v
	}

	var newJobName string
	if len(cronJob.Name) < 42 {
		newJobName = cronJob.Name + "-manual-" + rand.String(3)
	} else {
		newJobName = cronJob.Name[0:41] + "-manual-" + rand.String(3)
	}

	jobToCreate := &batch.Job{
		ObjectMeta: metaV1.ObjectMeta{
			Name:        newJobName,
			Namespace:   namespace,
			Annotations: annotations,
			Labels:      labels,
		},
		Spec: cronJob.Spec.JobTemplate.Spec,
	}

	_, err = client.BatchV1().Jobs(namespace).Create(context.TODO(), jobToCreate, metaV1.CreateOptions{})

	if err != nil {
		return err
	}
	return nil
}

func filterJobsByOwnerUID(UID types.UID, jobs []batch.Job) (matchingJobs []batch.Job) {
	for _, j := range jobs {
		for _, i := range j.OwnerReferences {
			if i.UID == UID {
				matchingJobs = append(matchingJobs, j)
				break
			}
		}
	}
	return
}

func filterJobsByState(active bool, jobs []batch.Job) (matchingJobs []batch.Job) {
	for _, j := range jobs {
		if active && j.Status.Active > 0 {
			matchingJobs = append(matchingJobs, j)
		} else if !active && j.Status.Active == 0 {
			matchingJobs = append(matchingJobs, j)
		}
	}
	return
}
