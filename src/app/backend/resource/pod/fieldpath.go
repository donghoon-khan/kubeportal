package pod

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/api/meta"
)

func FormatMap(m map[string]string) (fmtStr string) {
	for key, value := range m {
		fmtStr += fmt.Sprintf("%v=%q\n", key, value)
	}

	fmtStr = strings.TrimSuffix(fmtStr, "\n")
	return
}

func ExtractFieldPathAsString(obj interface{}, fieldPath string) (string, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return "", nil
	}

	switch fieldPath {
	case "metadata.annotations":
		return FormatMap(accessor.GetAnnotations()), nil
	case "metadata.labels":
		return FormatMap(accessor.GetLabels()), nil
	case "metadata.name":
		return accessor.GetName(), nil
	case "metadata.namespace":
		return accessor.GetNamespace(), nil
	}
	return "", fmt.Errorf("unsupported fieldPath: %v", fieldPath)
}
