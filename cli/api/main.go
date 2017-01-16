package main

import (
	"fmt"
	"github.com/AdamJacobMuller/home-api/api/server"
	"github.com/namsral/flag"
	"os"
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

	api := apiserver.NewAPIServer()
	api.LoadAndCreateProviders(configuration)
	api.Serve()
}
