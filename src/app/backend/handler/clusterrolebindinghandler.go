package handler

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"

	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	"github.com/donghoon-khan/kubeportal/src/app/backend/handler/parser"
	"github.com/donghoon-khan/kubeportal/src/app/backend/resource/clusterrolebinding"
)

func (apiHandler *APIHandler) installClusterRoleBinding(ws *restful.WebService) {
	ws.Route(
		ws.GET("/clusterrolebinding").
			To(apiHandler.handleGetClusterRoleBindingList).
			Writes(clusterrolebinding.ClusterRoleBindingList{}))
	ws.Route(
		ws.GET("/clusterrolebinding/{clusterrolebinding}").
			To(apiHandler.handleGetClusterRoleBindingDetail).
			Writes(clusterrolebinding.ClusterRoleBindingDetail{}))
}

// handleGetClusterRoleBindingList godoc
// @Tags ClusterRoleBinding
// @Summary Get list of clusterrolebinding
// @Description Returns a list of clusterrolebinding
// @Accept  json
// @Produce  json
// @Router /clusterrolebinding [GET]
// @Param itemsPerPage query int false "The number of items per page can be configured adding a query parameter named itemsPerPage"
// @Param page query int false "The number of page"
// @Param sortBy query string false "be used to sort the result list in ascending or descending {name, creationTimestamp, namespace, status, type} `e.g. sortBy=asc,name`"
// @Param filterBy query string false "be used to get filterd result list {name, creationTimestamp, namespace, status, type} `e.g. filterBy=namespace,kube-system`"
// @Success 200 {object} clusterrolebinding.ClusterRoleBindingList
// @Failure 401 {string} string "Unauthorized"
func (apiHandler *APIHandler) handleGetClusterRoleBindingList(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	dataSelect := parser.ParseDataSelectPathParameter(request)
	result, err := clusterrolebinding.GetClusterRoleBindingList(k8s, dataSelect)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

// handleGetClusterRoleBindingDetail godoc
// @Tags ClusterRoleBinding
// @Summary Get detail of clusterrolebinding
// @Description Returns a detail of clusterrolebinding
// @Accept  json
// @Produce  json
// @Router /clusterrolebinding/{clusterrolebinding} [GET]
// @Param clusterrolebinding path string true "Name of clusterrolebinding"
// @Success 200 {object} clusterrolebinding.ClusterRoleBindingDetail
// @Failure 401 {string} string "Unauthorized"
func (apiHandler *APIHandler) handleGetClusterRoleBindingDetail(request *restful.Request, response *restful.Response) {
	k8s, err := apiHandler.kManager.Kubernetes(request)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}

	crbName := request.PathParameter("clusterrolebinding")
	result, err := clusterrolebinding.GetClusterRoleBindingDetail(k8s, crbName)
	if err != nil {
		errors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}
