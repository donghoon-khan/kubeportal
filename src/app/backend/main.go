package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/swaggo/http-swagger"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"

	"github.com/donghoon-khan/kubeportal/src/app/backend/args"
	"github.com/donghoon-khan/kubeportal/src/app/backend/auth"
	authApi "github.com/donghoon-khan/kubeportal/src/app/backend/auth/api"
	"github.com/donghoon-khan/kubeportal/src/app/backend/handler"
	"github.com/donghoon-khan/kubeportal/src/app/backend/kubernetes"
	k8sApi "github.com/donghoon-khan/kubeportal/src/app/backend/kubernetes/api"
)

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

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:9090/swagger/doc.json"), //The url pointing to API definition"
	))

	r.Handle("/api/*", apiHandler)

	//http.Handle("/api/", apiHandler)
	//http.Handle("/", r)

	/*r := mux.NewRouter()
	r.Handle("/", apiHandler)
	http.Handle("/", r)*/

	//r := mux.NewRouter()
	//r.Handle("/api/", apiHandler)
	//mux.Handle("/api/", apiHandler)

	/*r := mux.NewRouter()
	r.Handle("/api/", apiHandler)

	http.Handle("/", r)*/

	/*r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:9090/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("#swagger-ui"),
	))*/

	//go func() { log.Fatal(http.ListenAndServe(":9090", nil)) }()
	go func() { log.Fatal(http.ListenAndServe(":9090", r)) }()
	select {}
}

func initAuthManager(k8sManager k8sApi.KubernetesManager) authApi.AuthManager {

	return auth.NewAuthManager(k8sManager,
		authApi.AuthenticationModes{authApi.Token: true, authApi.Basic: true},
		true)

}

func initArgHolder() {
	builder := args.GetHolderBuilder()
	builder.SetApiServerHost("https://127.0.0.1:16443")
	builder.SetApiLogLevel("INFO")
	builder.SetKubeConfigFile("/var/snap/microk8s/current/credentials/client.config")
	builder.SetNamespace("default")
	builder.SetPort(9090)
}

func handleFatalInitError(err error) {
	log.Fatalf("Error while initializing connection to Kubernetes apiserver. "+
		"This most likely means that the cluster is misconfigured (e.g., it has "+
		"invalid apiserver certificates or service account's configuration) or the "+
		"--apiserver-host param points to a server that does not exist. Reason: %s\n", err)
}
