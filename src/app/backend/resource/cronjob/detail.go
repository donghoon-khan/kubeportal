package cronjob

import (
	"context"

	batch2 "k8s.io/api/batch/v1beta1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sClient "k8s.io/client-go/kubernetes"
)

type CronJobDetail struct {
	CronJob                 `json:",inline"`
	ConcurrencyPolicy       string  `json:"concurrencyPolicy"`
	StartingDeadLineSeconds *int64  `json:"startingDeadlineSeconds"`
	Errors                  []error `json:"errors"`
}

func GetCronJobDetail(client k8sClient.Interface, namespace, name string) (*CronJobDetail, error) {
	rawObject, err := client.BatchV1beta1().CronJobs(namespace).Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	cj := toCronJobDetail(rawObject)
	return &cj, nil
}

func toCronJobDetail(cj *batch2.CronJob) CronJobDetail {
	return CronJobDetail{
		CronJob:                 toCronJob(cj),
		ConcurrencyPolicy:       string(cj.Spec.ConcurrencyPolicy),
		StartingDeadLineSeconds: cj.Spec.StartingDeadlineSeconds,
	}
}
