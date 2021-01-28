package cronjob

import (
	"log"

	"k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/donghoon-khan/kubeportal/src/app/backend/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	metricApi "github.com/donghoon-khan/kubeportal/src/app/backend/integration/metric/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
)

type CronJobList struct {
	ListMeta          api.ListMeta          `json:"listMeta"`
	CumulativeMetrics []metricApi.Metric    `json:"cumulativeMetrics"`
	Items             []CronJob             `json:"items"`
	Status            common.ResourceStatus `json:"status"`
	Errors            []error               `json:"errors"`
}

type CronJob struct {
	ObjectMeta   api.ObjectMeta `json:"objectMeta"`
	TypeMeta     api.TypeMeta   `json:"typeMeta"`
	Schedule     string         `json:"schedule"`
	Suspend      *bool          `json:"suspend"`
	Active       int            `json:"active"`
	LastSchedule *metav1.Time   `json:"lastSchedule"`
}

func GetCronJobList(kubernetes kubernetes.Interface, nsQuery *common.NamespaceQuery,
	dsQuery *dataselect.DataSelectQuery, metricClient metricApi.MetricClient) (*CronJobList, error) {
	log.Print("Getting list of all cron jobs in the cluster")

	channels := &common.ResourceChannels{
		CronJobList: common.GetCronJobListChannel(kubernetes, nsQuery, 1),
	}

	return GetCronJobListFromChannels(channels, dsQuery, metricClient)
}

func GetCronJobListFromChannels(channels *common.ResourceChannels, dsQuery *dataselect.DataSelectQuery,
	metricClient metricApi.MetricClient) (*CronJobList, error) {

	cronJobs := <-channels.CronJobList.List
	err := <-channels.CronJobList.Error
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	cronJobList := toCronJobList(cronJobs.Items, nonCriticalErrors, dsQuery, metricClient)
	cronJobList.Status = getStatus(cronJobs)
	return cronJobList, nil
}

func toCronJobList(cronJobs []v1beta1.CronJob, nonCriticalErrors []error, dsQuery *dataselect.DataSelectQuery,
	metricClient metricApi.MetricClient) *CronJobList {

	list := &CronJobList{
		Items:    make([]CronJob, 0),
		ListMeta: api.ListMeta{TotalItems: len(cronJobs)},
		Errors:   nonCriticalErrors,
	}

	cachedResources := &metricApi.CachedResources{}

	cronJobCells, metricPromises, filteredTotal := dataselect.GenericDataSelectWithFilterAndMetrics(ToCells(cronJobs),
		dsQuery, cachedResources, metricClient)
	cronJobs = FromCells(cronJobCells)
	list.ListMeta = api.ListMeta{TotalItems: filteredTotal}

	for _, cronJob := range cronJobs {
		list.Items = append(list.Items, toCronJob(&cronJob))
	}

	cumulativeMetrics, err := metricPromises.GetMetrics()
	if err != nil {
		list.CumulativeMetrics = make([]metricApi.Metric, 0)
	} else {
		list.CumulativeMetrics = cumulativeMetrics
	}

	return list
}

func toCronJob(cj *v1beta1.CronJob) CronJob {
	return CronJob{
		ObjectMeta:   api.NewObjectMeta(cj.ObjectMeta),
		TypeMeta:     api.NewTypeMeta(api.ResourceKindCronJob),
		Schedule:     cj.Spec.Schedule,
		Suspend:      cj.Spec.Suspend,
		Active:       len(cj.Status.Active),
		LastSchedule: cj.Status.LastScheduleTime,
	}
}
