package errors

import (
	"log"
	"net/http"

	"github.com/emicklei/go-restful"
	"k8s.io/apimachinery/pkg/api/errors"
)

var NonCriticalErrors = []int32{http.StatusForbidden, http.StatusUnauthorized}

func HandleError(err error) ([]error, error) {
	nonCriticalErrors := make([]error, 0)
	return AppendError(err, nonCriticalErrors)
}

func AppendError(err error, nonCriticalErrors []error) ([]error, error) {
	if err != nil {
		if isErrorCritical(err) {
			return nonCriticalErrors, LocalizeError(err)
		}
		log.Printf("Non-critical error occurred during resource retrieval: %s", err)
		nonCriticalErrors = appendMissing(nonCriticalErrors, LocalizeError(err))
	}
	return nonCriticalErrors, nil
}

func MergeErrors(errorArraysToMerge ...[]error) (mergedErrors []error) {
	for _, errorArray := range errorArraysToMerge {
		mergedErrors = appendMissing(mergedErrors, errorArray...)
	}
	return
}

func isErrorCritical(err error) bool {
	status, ok := err.(*errors.StatusError)
	if !ok {
		return true
	}
	return !contains(NonCriticalErrors, status.ErrStatus.Code)
}

func appendMissing(slice []error, toAppend ...error) []error {
	m := make(map[string]bool, 0)
	for _, s := range slice {
		m[s.Error()] = true
	}

	for _, a := range toAppend {
		_, ok := m[a.Error()]
		if !ok {
			slice = append(slice, a)
			m[a.Error()] = true
		}
	}

	return slice
}

func contains(s []int32, e int32) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func IsForbiddenError(err error) bool {
	status, ok := err.(*errors.StatusError)
	if !ok {
		return false
	}
	return status.ErrStatus.Code == http.StatusForbidden
}

func IsTokenExpiredError(err error) bool {
	if err == nil {
		return false
	}

	return err.Error() == MsgTokenExpiredError
}

func HandleInternalError(response *restful.Response, err error) {
	statusCode := http.StatusInternalServerError
	statusError, ok := err.(*errors.StatusError)
	if ok && statusError.Status().Code > 0 {
		statusCode = int(statusError.Status().Code)
	}
	response.AddHeader("Content-Type", "text/plain")
	response.WriteErrorString(statusCode, err.Error()+"\n")
}

func HandleHttpError(err error) int {
	if err == nil {
		return http.StatusInternalServerError
	}
	if err.Error() == MsgTokenExpiredError || err.Error() == MsgLoginUnauthorizedError || err.Error() == MsgEncryptionKeyChanged {
		return http.StatusUnauthorized
	}
	return http.StatusInternalServerError
}
