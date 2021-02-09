package handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/emicklei/go-restful/v3"

	"github.com/donghoon-khan/kubeportal/src/app/backend/args"
	"github.com/donghoon-khan/kubeportal/src/app/backend/errors"
	k8sapi "github.com/donghoon-khan/kubeportal/src/app/backend/kubernetes/api"
)

const (
	originalForwardedForHeader = "X-Original-Forwarded-For"
	forwardedForHeader         = "X-Forwarded-For"
	realIPHeader               = "X-Real-Ip"
)

func InstallFilters(ws *restful.WebService, manager k8sapi.KubernetesManager) {
	/*ws.Filter(requestAndResponseLogger)
	ws.Filter(metricsFilter)
	ws.Filter(validateXSRFFilter(manager.CSRFKey()))*/
	//ws.Filter(restrictedResourcesFilter)
}

func restrictedResourcesFilter(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {

	err := errors.NewUnauthorized(errors.MsgDashboardExclusiveResourceError)
	response.WriteHeaderAndEntity(int(err.ErrStatus.Code),
		errors.StatusErrorResponse{Message: err.Error()})
}

func requestAndResponseLogger(request *restful.Request, response *restful.Response,
	chain *restful.FilterChain) {
	if args.Holder.GetApiLogLevel() != "NONE" {
		log.Printf(formatRequestLog(request))
	}

	chain.ProcessFilter(request, response)

	if args.Holder.GetApiLogLevel() != "NONE" {
		log.Printf(formatResponseLog(response, request))
	}
}

func formatRequestLog(request *restful.Request) string {
	uri := ""
	content := "{}"

	if request.Request.URL != nil {
		uri = request.Request.URL.RequestURI()
	}

	byteArr, err := ioutil.ReadAll(request.Request.Body)
	if err == nil {
		content = string(byteArr)
	}

	request.Request.Body = ioutil.NopCloser(bytes.NewReader(byteArr))
	if args.Holder.GetApiLogLevel() != "DEBUG" && checkSensitiveURL(&uri) {
		content = "{ contents hidden }"
	}

	return fmt.Sprintf(RequestLogString, time.Now().Format(time.RFC3339), request.Request.Proto,
		request.Request.Method, uri, getRemoteAddr(request.Request), content)
}

func checkSensitiveURL(url *string) bool {
	var s struct{}
	var sensitiveUrls = make(map[string]struct{})
	sensitiveUrls["/api/v1/login"] = s
	sensitiveUrls["/api/v1/csrftoken/login"] = s
	sensitiveUrls["/api/v1/token/refresh"] = s

	if _, ok := sensitiveUrls[*url]; ok {
		return true
	}
	return false

}

// formatResponseLog formats response log string.
func formatResponseLog(response *restful.Response, request *restful.Request) string {
	return fmt.Sprintf(ResponseLogString, time.Now().Format(time.RFC3339),
		getRemoteAddr(request.Request), response.StatusCode())
}

func mapUrlToResource(url string) *string {
	parts := strings.Split(url, "/")
	if len(parts) < 3 {
		return nil
	}
	return &parts[3]
}

func getRemoteAddr(r *http.Request) string {
	if ip := getRemoteIPFromForwardHeader(r, originalForwardedForHeader); ip != "" {
		return ip
	}

	if ip := getRemoteIPFromForwardHeader(r, forwardedForHeader); ip != "" {
		return ip
	}

	if realIP := strings.TrimSpace(r.Header.Get(realIPHeader)); realIP != "" {
		return realIP
	}

	return r.RemoteAddr
}

func getRemoteIPFromForwardHeader(r *http.Request, header string) string {
	ips := strings.Split(r.Header.Get(header), ",")
	return strings.TrimSpace(ips[0])
}
