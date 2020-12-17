package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/pflag"

	"github.com/donghoon-khan/kubernetes-portal/src/app/backend/args"
	"github.com/donghoon-khan/kubernetes-portal/src/app/backend/handler"
)

var (
	argPort = pflag.Int("port", 9090, "The port to listen to for incoming HTTP requests.")
)

func initArgHolder() {
	builder := args.GetHolderBuilder()
	builder.SetPort(*argPort)
}

func main() {

	log.SetOutput(os.Stdout)

	initArgHolder()
	fmt.Println(args.Holder.GetPort())

	apiHandler, err := handler.CreateHTTPAPIHandler()
	if err != nil {
		handleFatalInitError(err)
	}

	server := &http.Server{
		Addr:         ":8080",
		Handler:      apiHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() { log.Fatal(server.ListenAndServe()) }()

	select {}
}

func handleFatalInitError(err error) {
	log.Fatalf("Error while initializing connection to Kubernetes apiserver. "+
		"This most likely means that the cluster is misconfigured (e.g., it has "+
		"invalid apiserver certificates or service account's configuration) or the "+
		"--apiserver-host param points to a server that does not exist. Reason: %s\n", err)
}
