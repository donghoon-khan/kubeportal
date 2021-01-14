package main

import (
	"flag"
	"log"
	"os"

	"github.com/donghoon-khan/kubeportal/src/app/backend/args"
	"github.com/spf13/pflag"
)

var (
	argPort          = pflag.Int("port", 3000, "The port to listen to for incoming HTTP requests.")
	argApiserverHost = pflag.String("apiserver-host", "", "The address of the Kubernetes Apiserver "+
		"to connect to in the format of protocol://address:port, e.g., "+
		"http://localhost:8080. If not specified, the assumption is that the binary runs inside a "+
		"Kubernetes cluster and local discovery is attempted.")
	argKubeConfigFile = pflag.String("kubeconfig", "", "Path to kubeconfig file with authorization and master location information.")
	argAPILogLevel    = pflag.String("api-log-level", "INFO", "Level of API request logging. Should be one of 'INFO|NONE|DEBUG'. Default: 'INFO'.")
)

func main() {
	log.SetOutput(os.Stdout)

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	flag.CommandLine.Parse(make([]string, 0))

	initArgHolder()

	log.Printf("Service Start")
	if args.Holder.GetApiServerHost() != "" {
		log.Printf("Using apiserver-host location: %s", args.Holder.GetApiServerHost())
	}
	if args.Holder.GetKubeConfigFile() != "" {
		log.Printf("Using kubeconfig file: %s", args.Holder.GetKubeConfigFile())
	}
}

func initArgHolder() {
	builder := args.GetHolderBuilder()
	builder.SetPort(*argPort)
	builder.SetApiServerHost(*argApiserverHost)
	builder.SetKubeConfigFile(*argKubeConfigFile)
	builder.SetApiLogLevel(*argAPILogLevel)
}
