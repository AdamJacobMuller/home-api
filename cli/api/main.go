package main

import (
	"github.com/AdamJacobMuller/home-api/api/server"
)

func main() {
	api := apiserver.NewAPIServer()

	api.LoadAndCreateProviders("config.json")

	api.Serve()
}
