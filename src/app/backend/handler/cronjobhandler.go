package handler

import (
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"

	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/handler/parser"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/common"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/cronjob"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/dataselect"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/job"
)

func (apiHandler *APIHandler) installCronJob(ws *restful.WebService) {
	ws.Route(
		ws.GET("/cronjob").
			To(apiHandler.handleGetCronJobList).
			Returns(200, "OK", cronjob.CronJobList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List objects of kind CronJob").
			Metadata(restfulspec.KeyOpenAPITags, []string{cronJobDocsTag}))
	ws.Route(
		ws.GET("/cronjob/{namespace}").
			To(apiHandler.handleGetCronJobListNamespace).
			Param(ws.PathParameter("namespace", "Query for Namespace").Required(true)).
			Returns(200, "OK", cronjob.CronJobList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List objects of kind CronJob in the Namespace").
			Metadata(restfulspec.KeyOpenAPITags, []string{cronJobDocsTag}))
	ws.Route(
		ws.GET("/cronjob/{namespace}/{name}").
			To(apiHandler.handleGetCronJobDetail).
			Param(ws.PathParameter("namespace", "Query for Namespace").Required(true)).
			Param(ws.PathParameter("name", "Name of CronJob").DataType("string").Required(true)).
			Returns(200, "OK", cronjob.CronJobDetail{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("Read the specified CronJob").
			Metadata(restfulspec.KeyOpenAPITags, []string{cronJobDocsTag}))
	ws.Route(
		ws.GET("/cronjob/{namespace}/{name}/job").
			To(apiHandler.handleGetCronJobJobs).
			Param(ws.PathParameter("namespace", "Query for Namespace").Required(true)).
			Param(ws.PathParameter("name", "Name of CronJob").DataType("string").Required(true)).
			Returns(200, "OK", job.JobList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List Jobs related to a CronJob").
			Metadata(restfulspec.KeyOpenAPITags, []string{cronJobDocsTag}))
	ws.Route(
		ws.GET("/cronjob/{namespace}/{name}/event").
			To(apiHandler.handleGetCronJobEvents).
			Param(ws.PathParameter("namespace", "Query for Namespace").Required(true)).
			Param(ws.PathParameter("name", "Name of CronJob").DataType("string").Required(true)).
			Returns(200, "OK", common.EventList{}).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("List events related to a CronJob").
			Metadata(restfulspec.KeyOpenAPITags, []string{cronJobDocsTag}))
	ws.Route(
		ws.PUT("/cronjob/{namespace}/{name}/trigger").
			To(apiHandler.handleTriggerCronJob).
			Param(ws.PathParameter("namespace", "Query for Namespace").Required(true)).
			Param(ws.PathParameter("name", "Name of CronJob").DataType("string").Required(true)).
			Returns(200, "OK", nil).
			Returns(401, "Unauthorized", errors.StatusErrorResponse{}).
			Doc("Replace trigger related to a CronJob").
			Metadata(restfulspec.KeyOpenAPITags, []string{cronJobDocsTag}))
}

func (apiHandler *APIHandler) handleGetCronJobList(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	dataSelect := parser.ParseDataSelectPathParameter(request)
	dataSelect.MetricQuery = dataselect.StandardMetrics
	result, err := cronjob.GetCronJobList(k8s, common.NewNamespaceQuery(nil), dataSelect, apiHandler.iManager.Metric().Client())
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetCronJobListNamespace(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := parseNamespacePathParameter(request)
	dataSelect := parser.ParseDataSelectPathParameter(request)
	dataSelect.MetricQuery = dataselect.StandardMetrics
	result, err := cronjob.GetCronJobList(k8s, namespace, dataSelect, apiHandler.iManager.Metric().Client())
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetCronJobDetail(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")
	result, err := cronjob.GetCronJobDetail(k8s, namespace, name)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetCronJobJobs(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")
	active := true
	if request.QueryParameter("active") == "false" {
		active = false
	}

	dataSelect := parser.ParseDataSelectPathParameter(request)
	result, err := cronjob.GetCronJobJobs(k8s, apiHandler.iManager.Metric().Client(), dataSelect, namespace, name, active)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetCronJobEvents(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")
	dataSelect := parser.ParseDataSelectPathParameter(request)
	result, err := cronjob.GetCronJobEvents(k8s, dataSelect, namespace, name)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleTriggerCronJob(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")
	err = cronjob.TriggerCronJob(k8s, namespace, name)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeader(http.StatusOK)
}
