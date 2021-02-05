package handler

import (
	"net/http"

	"github.com/donghoon-khan/kubeportal/src/app/backend/docs"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/spec"
)

const (
	authenticationDocsTag        = "Authentication"
	integrationDocsTag           = "Integration"
	clusterRoleBindingDocsTag    = "ClusterRoleBinding"
	clusterRoleDocsTag           = "ClusterRole"
	configMapDocsTag             = "ConfigMap"
	cronJobDocsTag               = "CronJob"
	ingressDocsTag               = "Ingress"
	nodeDocsTag                  = "Node"
	persistentVolumeClaimDocsTag = "PersistentVolumeClaim"
	podDocsTag                   = "Pod"
	secretDocsTag                = "Sceret"
	serviceDocsTag               = "Service"
	serviceAccountDocsTag        = "ServiceAccount"
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
				Name:        docs.AuthenticationDocsTag,
				Description: "Before kubernetes-portal API, You must be authenticated using Authentication API.",
			},
		},
		{
			TagProps: spec.TagProps{
				Name: docs.IntegrationDocsTag,
				Description: "Currently Dashboard implements metrics-server and Heapster integrations." +
					" They are using integration framework that allows to support and integrate more metric providers as well as additional applications such as Weave Scope or Grafana." +
					"<br/>Ref: https://github.com/kubernetes-sigs/metrics-server or https://github.com/kubernetes-retired/heapster",
			},
		},
		{
			TagProps: spec.TagProps{
				Name: docs.ClusterRoleBindingDocsTag,
				Description: "ClusterRoleBinding references a ClusterRole, but not contain it." +
					"It can reference a ClusterRole in the global namespace, and adds who information via Subject." +
					"<br/>Ref: https://kubernetes.io/docs/reference/access-authn-authz/rbac/#rolebinding-and-clusterrolebinding",
			},
		},
		{
			TagProps: spec.TagProps{
				Name: docs.ClusterRoleDocsTag,
				Description: "ClusterRole is a logical grouping of PolicyRules that can be referenced as a unit by ClusterRoleBindings." +
					"<br/>Ref: https://kubernetes.io/docs/reference/access-authn-authz/rbac/#role-and-clusterrole",
			},
		},
		{
			TagProps: spec.TagProps{
				Name: docs.ConfigMapDocsTag,
				Description: "A ConfigMap is an API object used to store non-confidential data in key-value pairs. Pods can consume ConfigMaps as environment variables, command-line arguments, or as configuration files in a volume." +
					"<br/>Ref: https://kubernetes.io/docs/concepts/configuration/configmap/",
			},
		},
		{
			TagProps: spec.TagProps{
				Name: docs.CronJobDocsTag,
				Description: "One CronJob object is like one line of a crontab (cron table) file. It runs a job periodically on a given schedule, written in [Cron](https://en.wikipedia.org/wiki/Cron) format." +
					"<br/>Ref: https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/",
			},
		},
		{
			TagProps: spec.TagProps{
				Name: docs.IngressDocsTag,
				Description: "Ingress is a collection of rules that allow inbound connections to reach the endpoints defined by a backend. An Ingress can be configured to give services externally-reachable urls, load balance traffic, terminate SSL, offer name based virtual hosting etc." +
					"<br/>Ref: https://kubernetes.io/docs/concepts/services-networking/ingress/",
			},
		},
		{
			TagProps: spec.TagProps{
				Name: docs.NodeDocsTag,
				Description: "Node is a worker node in Kubernetes. Each node will have a unique identifier in the cache (i.e. in etcd)." +
					"<br/>Ref: https://kubernetes.io/docs/concepts/architecture/nodes/",
			},
		},
		{
			TagProps: spec.TagProps{
				Name: docs.PersistentVolumeClaimDocsTag,
				Description: "PersistentVolumeClaim is a user's request for and claim to a persistent volume" +
					"<br/>Ref: https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims",
			},
		},
		{
			TagProps: spec.TagProps{
				Name: docs.PodDocsTag,
				Description: "Pod is a collection of containers that can run on a host. This resource is created by clients and scheduled onto hosts." +
					"<br/>Ref: https://kubernetes.io/docs/concepts/workloads/pods/",
			},
		},
		{
			TagProps: spec.TagProps{
				Name: docs.SecretDocsTag,
				Description: "Secret holds secret data of a certain type. The total bytes of the values in the Data field must be less than MaxSecretSize bytes." +
					"<br/>Ref: https://kubernetes.io/docs/concepts/configuration/secret/",
			},
		},
		{
			TagProps: spec.TagProps{
				Name: docs.ServiceDocsTag,
				Description: "Service is a named abstraction of software service (for example, mysql) consisting of local port (for example 3306) that the proxy listens on, and the selector that determines which pods will answer requests sent through the proxy." +
					"<br/>Ref: https://kubernetes.io/docs/concepts/services-networking/service/",
			},
		},
		{
			TagProps: spec.TagProps{
				Name: docs.ServiceAccountDocsTag,
				Description: "ServiceAccount binds together: * a name, understood by users, and perhaps by peripheral systems, for an identity * a principal that can be authenticated and authorized * a set of secrets" +
					"<br/>Ref: https://kubernetes.io/docs/reference/access-authn-authz/service-accounts-admin/",
			},
		},
	}
}
