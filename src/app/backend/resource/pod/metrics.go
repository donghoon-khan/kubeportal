package pod

/*
import (
	"log"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	metricApi "github.com/donghoon-khan/kubeportal/src/app/backend/integration/metric/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
	"github.com/kubernetes/dashboard/src/app/backend/api"
)

type MetricsByPod struct {
	MetricsMap map[types.UID]PodMetrics `json:"metricsMap"`
}

type PodMetrics struct {
	CPUUsage           *uint64                 `json:"cpuUsage"`
	MemoryUsage        *uint64                 `json:"memoryUsage"`
	CPUUsageHistory    []metricApi.MetricPoint `json:"cpuUsageHistory"`
	MemoryUsageHistory []metricApi.MetricPoint `json:"memoryUsageHistory"`
}

func getMetricsPerPod(pods []v1.Pod, metricClient metricApi.MetricClient, dsQuery *dataselect.DataSelectQuery) (
	*MetricsByPod, error) {
	log.Println("Getting pod metrics")

	result := &MetricsByPod{MetricsMap: make(map[types.UID]PodMetrics)}

	metricPromises := dataselect.PodListMetrics(toCells(pods), dsQuery, metricClient)
	metrics, err := metricPromises.GetMetrics()
	if err != nil {
		return result, err
	}

	for _, m := range metrics {
		uid, err := getPodUIDFromMetric(m)
		if err != nil {
			log.Printf("Skipping metric because of error: %s", err.Error())
		}

		podMetrics := PodMetrics{}
		if p, exists := result.MetricsMap[uid]; exists {
			podMetrics = p
		}

		if m.MetricName == metricApi.CpuUsage && len(m.MetricPoints) > 0 {
			podMetrics.CPUUsage = &m.MetricPoints[len(m.MetricPoints)-1].Value
			podMetrics.CPUUsageHistory = m.MetricPoints
		}

		if m.MetricName == metricApi.MemoryUsage && len(m.MetricPoints) > 0 {
			podMetrics.MemoryUsage = &m.MetricPoints[len(m.MetricPoints)-1].Value
			podMetrics.MemoryUsageHistory = m.MetricPoints
		}

		result.MetricsMap[uid] = podMetrics
	}

	return result, nil
}

func getPodUIDFromMetric(metric metricApi.Metric) (types.UID, error) {
	uidList, exists := metric.Label[api.ResourceKindPod]
	if !exists {
		return "", errors.NewInvalid("Metric label not set.")
	}

	if len(uidList) != 1 {
		return "", errors.NewInvalid("Found multiple UIDs. Metric should contain data for single resource only.")
	}

	return uidList[0], nil
}
*/
