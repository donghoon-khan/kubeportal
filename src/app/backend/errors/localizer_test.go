package errors_test

import (
	"testing"

	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
)

func TestLocalizeError(t *testing.T) {
	cases := []struct {
		err      error
		expected error
	}{
		{
			nil,
			nil,
		},
		{
			errors.NewInternal("some unknown error"),
			errors.NewInternal("some unknown error"),
		},
		{
			errors.NewInvalid("does not match the namespace"),
			errors.NewInvalid("MSG_DEPLOY_NAMESPACE_MISMATCH_ERROR"),
		},
		{
			errors.NewInvalid("empty namespace may not be set"),
			errors.NewInvalid("MSG_DEPLOY_EMPTY_NAMESPACE_ERROR"),
		},
		{
			errors.NewInvalid("the server has asked for the client to provide credentials"),
			errors.NewInvalid("MSG_LOGIN_UNAUTHORIZED_ERROR"),
		},
	}

	for _, c := range cases {
		actual := errors.LocalizeError(c.err)
		if !areErrorsEqual(actual, c.expected) {
			t.Errorf("LocalizeError(%+v) == %+v, expected %+v", c.err, actual, c.expected)
		}
	}
}

func areErrorsEqual(err1, err2 error) bool {
	return (err1 != nil && err2 != nil && err1.Error() == err2.Error() ||
		(err1 == nil && err2 == nil))
}
