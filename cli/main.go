package main

import (
	"github.com/AdamJacobMuller/home-api/api/controller"
	"github.com/AdamJacobMuller/home-api/api/server"
	"github.com/AdamJacobMuller/home-api/providers/homeseer"
)

func main() {
	controller := apicontroller.NewAPIController()
	hscontroller := homeseer.NewHomeseerController("http://homeseer.adam.gs")
	api := apiserver.NewAPIServer()

	controller.AddProvider(hscontroller)
	api.Controller = controller

	api.Serve()
}
