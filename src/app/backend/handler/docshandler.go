package handler

import (
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/spec"
)

func CreateApiDocsHTTPHandler(wsContainer *restful.Container, specURL string, next http.Handler) http.Handler {

	config := restfulspec.Config{
		WebServices:                   wsContainer.RegisteredWebServices(),
		APIPath:                       "/apidocs.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject}

	restful.DefaultContainer.Add(restfulspec.NewOpenAPIService(config))

	opts := middleware.RedocOpts{SpecURL: specURL}
	sh := middleware.Redoc(opts, next)

	return sh
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "Fission OpenAPI 2.0",
			Description: "TEST",
			Version:     "v1",
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{Name: "dhkang"},
			},
		},
	}
	swo.Tags = []spec.Tag{
		{
			TagProps: spec.TagProps{
				Name:        "Authentication",
				Description: "Before kubernetes-portal API, You must be authenticated using Authentication API.",
			},
		},
		{
			TagProps: spec.TagProps{
				Name: "Integration",
				Description: "Currently Dashboard implements metrics-server and Heapster integrations." +
					" They are using integration framework that allows to support and integrate more metric providers as well as additional applications such as Weave Scope or Grafana." +
					"<br/>Ref: https://github.com/kubernetes-sigs/metrics-server or https://github.com/kubernetes-retired/heapster",
			},
		},
		{
			TagProps: spec.TagProps{
				Name: "ClusterRoleBinding",
				Description: "ClusterRoleBinding references a ClusterRole, but not contain it." +
					"It can reference a ClusterRole in the global namespace, and adds who information via Subject." +
					"<br/>Ref: https://kubernetes.io/docs/reference/access-authn-authz/rbac/#rolebinding-and-clusterrolebinding",
			},
		},
		{
			TagProps: spec.TagProps{
				Name: "ClusterRole",
				Description: "ClusterRole is a logical grouping of PolicyRules that can be referenced as a unit by ClusterRoleBindings." +
					"<br/>Ref: https://kubernetes.io/docs/reference/access-authn-authz/rbac/#role-and-clusterrole",
			},
		},
		{
			TagProps: spec.TagProps{
				Name: "Node",
				Description: "Node is a worker node in Kubernetes. Each node will have a unique identifier in the cache (i.e. in etcd)." +
					"<br/>Ref: https://kubernetes.io/docs/concepts/architecture/nodes/",
			},
		},
		{
			TagProps: spec.TagProps{
				Name: "PersistentVolumeClaim",
				Description: "PersistentVolumeClaim is a user's request for and claim to a persistent volume" +
					"<br/>Ref: https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims",
			},
		},
		{
			TagProps: spec.TagProps{
				Name: "Pod",
				Description: "Pod is a collection of containers that can run on a host. This resource is created by clients and scheduled onto hosts." +
					"<br/>Ref: https://kubernetes.io/docs/concepts/workloads/pods/",
			},
		},
		{
			TagProps: spec.TagProps{
				Name: "Secret",
				Description: "Secret holds secret data of a certain type. The total bytes of the values in the Data field must be less than MaxSecretSize bytes." +
					"<br/>Ref: https://kubernetes.io/docs/concepts/configuration/secret/",
			},
		}}
}
