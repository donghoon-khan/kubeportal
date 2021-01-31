package container

import (
	"context"
	"io"
	"io/ioutil"

	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/logs"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

var lineReadLimit int64 = 5000

var byteReadLimit int64 = 500000

type PodContainerList struct {
	Containers []string `json:"containers"`
}

func GetPodContainers(kubernetes kubernetes.Interface, ns, podID string) (*PodContainerList, error) {
	pod, err := kubernetes.CoreV1().Pods(ns).Get(context.TODO(), podID, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	containers := &PodContainerList{Containers: make([]string, 0)}
	for _, container := range pod.Spec.Containers {
		containers.Containers = append(containers.Containers, container.Name)
	}
	return containers, nil
}

func GetLogDetail(kubernetes kubernetes.Interface, ns, podID, container string,
	logSelector *logs.Selection, usePreviousLogs bool) (*logs.LogDetails, error) {
	pod, err := kubernetes.CoreV1().Pods(ns).Get(context.TODO(), podID, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if len(container) == 0 {
		container = pod.Spec.Containers[0].Name
	}
	logOptions := mapToLogOptions(container, logSelector, usePreviousLogs)
	rawLogs, err := readRawLogs(kubernetes, ns, podID, logOptions)
	if err != nil {
		return nil, err
	}
	details := ConstructLogDetails(podID, rawLogs, container, logSelector)
	return details, nil
}

func GetLogFile(kubernetes kubernetes.Interface, ns, podID, container string, usePreviousLogs bool) (io.ReadCloser, error) {
	logOptions := &v1.PodLogOptions{
		Container:  container,
		Follow:     false,
		Previous:   usePreviousLogs,
		Timestamps: false,
	}
	logStream, err := openStream(kubernetes, ns, podID, logOptions)
	return logStream, err
}

func ConstructLogDetails(podID, rawLogs, container string, logSelector *logs.Selection) *logs.LogDetails {
	parsedLines := logs.ToLogLines(rawLogs)
	logLines, fromDate, toDate, logSelection, lastPage := parsedLines.SelectLogs(logSelector)

	readLimitReached := isReadLimitReached(int64(len(rawLogs)), int64(len(parsedLines)), logSelector.LogFilePosition)
	truncated := readLimitReached && lastPage

	info := logs.LogInfo{
		PodName:       podID,
		ContainerName: container,
		FromDate:      fromDate,
		ToDate:        toDate,
		Truncated:     truncated,
	}
	return &logs.LogDetails{
		Info:      info,
		Selection: logSelection,
		LogLines:  logLines,
	}
}

func mapToLogOptions(container string, logSelector *logs.Selection, previous bool) *v1.PodLogOptions {
	logOptions := &v1.PodLogOptions{
		Container:  container,
		Follow:     false,
		Previous:   previous,
		Timestamps: true,
	}

	if logSelector.LogFilePosition == logs.Beginning {
		logOptions.LimitBytes = &byteReadLimit
	} else {
		logOptions.TailLines = &lineReadLimit
	}

	return logOptions
}

func readRawLogs(kubernetes kubernetes.Interface, ns, podID string, logOptions *v1.PodLogOptions) (
	string, error) {
	readCloser, err := openStream(kubernetes, ns, podID, logOptions)
	if err != nil {
		return err.Error(), nil
	}

	defer readCloser.Close()

	result, err := ioutil.ReadAll(readCloser)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func openStream(kubernetes kubernetes.Interface, ns, podID string, logOptions *v1.PodLogOptions) (io.ReadCloser, error) {
	return kubernetes.CoreV1().RESTClient().Get().
		Namespace(ns).
		Name(podID).
		Resource("pods").
		SubResource("log").
		VersionedParams(logOptions, scheme.ParameterCodec).Stream(context.TODO())
}

func isReadLimitReached(bytesLoaded int64, linesLoaded int64, logFilePosition string) bool {
	return (logFilePosition == logs.Beginning && bytesLoaded >= byteReadLimit) ||
		(logFilePosition == logs.End && linesLoaded >= lineReadLimit)
}
