package job

import (
	"context"

	batch "k8s.io/api/batch/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
)

type JobDetail struct {
	Job         `json:",inline"`
	Completions *int32  `json:"completions"`
	Errors      []error `json:"errors"`
}

func GetJobDetail(kubernetes kubernetes.Interface, namespace, name string) (*JobDetail, error) {
	jobData, err := kubernetes.BatchV1().Jobs(namespace).Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	podInfo, err := getJobPodInfo(kubernetes, jobData)
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	job := toJobDetail(jobData, *podInfo, nonCriticalErrors)
	return &job, nil
}

func toJobDetail(job *batch.Job, podInfo common.PodInfo, nonCriticalErrors []error) JobDetail {
	return JobDetail{
		Job:         toJob(job, &podInfo),
		Completions: job.Spec.Completions,
		Errors:      nonCriticalErrors,
	}
}
