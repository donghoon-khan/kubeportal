package errors_test

import (
	"reflect"
	"testing"

	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
)

func TestHandlerHttpError(t *testing.T) {
	cases := []struct {
		err      error
		expected int
	}{
		{
			nil,
			500,
		},
		{
			errors.NewInvalid("some unknown error"),
			500,
		},
		{
			errors.NewInvalid(errors.MsgDeployNamespaceMismatchError),
			500,
		},
		{
			errors.NewInvalid(errors.MsgDeployEmptyNamespaceError),
			500,
		},
		{
			errors.NewInvalid(errors.MsgLoginUnauthorizedError),
			401,
		},
		{
			errors.NewInvalid(errors.MsgEncryptionKeyChanged),
			401,
		},
		{
			errors.NewInvalid(errors.MsgDashboardExclusiveResourceError),
			500,
		},
		{
			errors.NewInvalid(errors.MsgTokenExpiredError),
			401,
		},
	}
	for _, c := range cases {
		actual := errors.HandleHttpError(c.err)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("HandleHttpError(%+v) == %+v, expected %+v", c.err, actual, c.expected)
		}
	}
}
