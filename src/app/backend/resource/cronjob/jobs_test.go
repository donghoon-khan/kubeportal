package cronjob_test

import (
	"context"
	"strings"
	"testing"

	batch "k8s.io/api/batch/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/cronjob"
)

func TestTriggerCronJobWithInvalidName(t *testing.T) {
	client := fake.NewSimpleClientset()

	err := cronjob.TriggerCronJob(client, namespace, "invalidName")
	if !errors.IsNotFound(err) {
		t.Error("TriggerCronJob should return error when invalid name is passed")
	}
}

func TestTriggerCronJobWithLongName(t *testing.T) {
	longName := strings.Repeat("test", 13)

	cron := batch.CronJob{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      longName,
			Namespace: namespace,
			Labels:    labels,
		}, TypeMeta: metaV1.TypeMeta{
			Kind:       "CronJob",
			APIVersion: "v1",
		}}

	client := fake.NewSimpleClientset(&cron)
	err := cronjob.TriggerCronJob(client, namespace, longName)
	if err != nil {
		t.Error(err)
	}
}

func TestTriggerCronJob(t *testing.T) {

	cron := batch.CronJob{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		}, TypeMeta: metaV1.TypeMeta{
			Kind:       "CronJob",
			APIVersion: "v1",
		}, Spec: batch.CronJobSpec{
			Schedule: "* * * * *",
			JobTemplate: batch.JobTemplateSpec{
				ObjectMeta: metaV1.ObjectMeta{
					Namespace: namespace,
					Labels:    labels,
				},
			},
		},
	}

	client := fake.NewSimpleClientset(&cron)

	err := cronjob.TriggerCronJob(client, namespace, name)
	if err != nil {
		t.Error(err)
	}

	list, err := client.BatchV1().Jobs(namespace).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		t.Error(err)
	}
	if len(list.Items) != 1 {
		t.Error(err)
	}
}
