package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-openapi/runtime/middleware"

	"github.com/donghoon-khan/kubeportal/src/app/backend/args"
	"github.com/donghoon-khan/kubeportal/src/app/backend/auth"

	authApi "github.com/donghoon-khan/kubeportal/src/app/backend/auth/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/handler"
	"github.com/donghoon-khan/kubeportal/src/app/backend/kubernetes"
	k8sApi "github.com/donghoon-khan/kubeportal/src/app/backend/kubernetes/api"
)

// @title Kubernetes-portal API
// @version 0.0.1
// @description This is a Kubernetes-portal api server
// @host localhost:9090
// @BasePath /api/v1
func main() {
	log.SetOutput(os.Stdout)
	initArgHolder()

	if args.Holder.GetApiServerHost() != "" {
		log.Printf("Using apiserver-host: %s", args.Holder.GetApiServerHost())
	}
	if args.Holder.GetKubeConfigFile() != "" {
		log.Printf("Using kubeconfig file: %s", args.Holder.GetKubeConfigFile())
	}
	if args.Holder.GetNamespace() != "" {
		log.Printf("Using namespace: %s", args.Holder.GetNamespace())
	}
	k8sManager := kubernetes.NewKubernetesManager(args.Holder.GetKubeConfigFile(), args.Holder.GetApiServerHost())
	versionInfo, err := k8sManager.InsecureKubernetes().Discovery().ServerVersion()
	if err != nil {
		handleFatalInitError(err)
	}
	log.Printf("Successful initial request to the apiserver, version: %s", versionInfo.String())

	authManager := initAuthManager(k8sManager)
	apiHandler, err := handler.CreateHttpApiHandler(k8sManager, authManager)
	if err != nil {
		handleFatalInitError(err)
	}

	r := chi.NewRouter()

	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)
	r.Handle("/docs", sh)
	r.Handle("/swagger.yaml", http.FileServer(http.Dir("./swagger/")))
	//r.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
	/*r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:9090/swagger/doc.json"), //The url pointing to API definition"
	))*/
	r.Handle("/api/*", apiHandler)

	go func() { log.Fatal(http.ListenAndServe(":9090", r)) }()
	select {}
}

func aamain() {
	//http.Handle("/", http.FileServer(http.Dir("./public")))
	//http.Handle("/static", http.FileServer(http.Dir("wwwroot")))
	//http.ListenAndServe(":5000", nil)
	http.ListenAndServe(":9090", http.FileServer(http.Dir("./swagger")))
}

func initAuthManager(k8sManager k8sApi.KubernetesManager) authApi.AuthManager {

	return auth.NewAuthManager(k8sManager,
		authApi.AuthenticationModes{authApi.Token: true, authApi.Basic: true},
		true)

}

func initArgHolder() {
	builder := args.GetHolderBuilder()
	builder.SetApiServerHost("http://127.0.0.1:8001")
	builder.SetApiLogLevel("INFO")
	builder.SetKubeConfigFile("/Users/kangdonghoon/.kube/config")
	//builder.SetKubeConfigFile("/var/snap/microk8s/current/credentials/client.config")
	builder.SetNamespace("default")
	builder.SetPort(9090)
}

func handleFatalInitError(err error) {
	log.Fatalf("Error while initializing connection to Kubernetes apiserver. "+
		"This most likely means that the cluster is misconfigured (e.g., it has "+
		"invalid apiserver certificates or service account's configuration) or the "+
		"--apiserver-host param points to a server that does not exist. Reason: %s\n", err)
}
