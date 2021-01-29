package pod

import (
	"strings"
	"testing"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestExtractFieldPathAsString(t *testing.T) {
	cases := []struct {
		name                    string
		fieldPath               string
		obj                     interface{}
		expectedValue           string
		expectedMessageFragment string
	}{
		{
			name:                    "not an API object",
			fieldPath:               "metadata.name",
			obj:                     "",
			expectedMessageFragment: "expected struct",
		},
		{
			name:      "ok - namespace",
			fieldPath: "metadata.namespace",
			obj: &v1.Pod{
				ObjectMeta: metaV1.ObjectMeta{
					Namespace: "object-namespace",
				},
			},
			expectedValue: "object-namespace",
		},
		{
			name:      "ok - name",
			fieldPath: "metadata.name",
			obj: &v1.Pod{
				ObjectMeta: metaV1.ObjectMeta{
					Name: "object-name",
				},
			},
			expectedValue: "object-name",
		},
		{
			name:      "ok - labels",
			fieldPath: "metadata.labels",
			obj: &v1.Pod{
				ObjectMeta: metaV1.ObjectMeta{
					Labels: map[string]string{"key": "value"},
				},
			},
			expectedValue: "key=\"value\"",
		},
		{
			name:      "ok - labels bslash n",
			fieldPath: "metadata.labels",
			obj: &v1.Pod{
				ObjectMeta: metaV1.ObjectMeta{
					Labels: map[string]string{"key": "value\n"},
				},
			},
			expectedValue: "key=\"value\\n\"",
		},
		{
			name:      "ok - annotations",
			fieldPath: "metadata.annotations",
			obj: &v1.Pod{
				ObjectMeta: metaV1.ObjectMeta{
					Annotations: map[string]string{"builder": "john-doe"},
				},
			},
			expectedValue: "builder=\"john-doe\"",
		},

		{
			name:      "invalid expression",
			fieldPath: "metadata.whoops",
			obj: &v1.Pod{
				ObjectMeta: metaV1.ObjectMeta{
					Namespace: "object-namespace",
				},
			},
			expectedMessageFragment: "unsupported fieldPath",
		},
	}

	for _, tc := range cases {
		actual, err := ExtractFieldPathAsString(tc.obj, tc.fieldPath)
		if err != nil {
			if tc.expectedMessageFragment != "" {
				if !strings.Contains(err.Error(), tc.expectedMessageFragment) {
					t.Errorf("%v: unexpected error message: %q, expected to contain %q", tc.name, err, tc.expectedMessageFragment)
				}
			} else {
				t.Errorf("%v: unexpected error: %v", tc.name, err)
			}
		} else if e := tc.expectedValue; e != "" && e != actual {
			t.Errorf("%v: unexpected result; got %q, expected %q", tc.name, actual, e)
		}
	}
}
