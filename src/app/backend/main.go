package main

import (
	"log"
	"net/http"
	"os"

	"github.com/donghoon-khan/kubeportal/src/app/backend/args"
	"github.com/donghoon-khan/kubeportal/src/app/backend/auth"
	"github.com/donghoon-khan/kubeportal/src/app/backend/docs"
	"github.com/donghoon-khan/kubeportal/src/app/backend/integration"

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
	iManager := integration.NewIntegrationManager(k8sManager)

	apiHandler, err := handler.CreateHttpApiHandler(iManager, k8sManager, authManager)
	if err != nil {
		handleFatalInitError(err)
	}

	http.Handle("/api/", apiHandler)
	http.Handle("/docs", docs.CreateApiDocsHTTPHandler(apiHandler, "/apidocs.json", nil))

	go func() { log.Fatal(http.ListenAndServe(":9090", nil)) }()
	select {}
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
	builder.SetKubeConfigFile("kube.config")
	builder.SetNamespace("default")
	builder.SetPort(9090)
}

func handleFatalInitError(err error) {
	log.Fatalf("Error while initializing connection to Kubernetes apiserver. "+
		"This most likely means that the cluster is misconfigured (e.g., it has "+
		"invalid apiserver certificates or service account's configuration) or the "+
		"--apiserver-host param points to a server that does not exist. Reason: %s\n", err)
}
