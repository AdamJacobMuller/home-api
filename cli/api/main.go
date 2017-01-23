package main

import (
	"fmt"
	"github.com/AdamJacobMuller/home-api/api/server"
	log "github.com/Sirupsen/logrus"
	"github.com/namsral/flag"
	"os"
)

var (
	Version   string
	BuildTime string
)

func main() {
	var configuration string
	var providers_d string

	flag.StringVar(&configuration, "configuration", "", "configuration file")
	flag.StringVar(&providers_d, "providers-d", "", "directory containing provider json configuration files")
	flag.Parse()

	if configuration == "" {
		fmt.Printf("-configuration is required\n")
		os.Exit(1)
	}
	log.WithFields(log.Fields{"version": Version, "build-time": BuildTime}).Info("starting up API Server")

	api := apiserver.NewAPIServer()
	api.LoadAndCreateProviders(configuration)
	api.Serve()
}
