package cronjob_test

import "github.com/donghoon-khan/kubeportal/src/app/backend/errors"

var (
	name         = "test-name"
	namespace    = "test-namespace"
	suspend      = false
	labels       = map[string]string{"app": "test-app"}
	eventMessage = "test-message"
	customError  = errors.NewInvalid("test-error")
)
