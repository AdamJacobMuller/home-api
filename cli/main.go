package main

import (
	"github.com/AdamJacobMuller/home-api/api"
	"github.com/AdamJacobMuller/home-api/providers/homeseer"
)

func main() {
	controller := homeseer.NewHomeseerController("http://homeseer.adam.gs")
	api := api.NewAPIServer()
	api.HSController = controller

	api.Serve()
	/*
		controller.SetDevicesValue(homeseer.HSLookup{"LocationOne": "Living Room", "TypeString": "Z-Wave Switch Binary"}, 0)
		controller.SetChildDevicesValue(homeseer.HSLookup{"LocationOne": "Living Room", "Child": homeseer.HSLookup{"Name": "Color Control Warm_White Channel"}}, 255)
		controller.SetChildDevicesValue(homeseer.HSLookup{"LocationOne": "Living Room", "Child": homeseer.HSLookup{"Name": "Color Control Cold_White Channel"}}, 0)
		controller.SetChildDevicesValue(homeseer.HSLookup{"LocationOne": "Living Room", "Child": homeseer.HSLookup{"Name": "Color Control Red Channel"}}, 0)
		controller.SetChildDevicesValue(homeseer.HSLookup{"LocationOne": "Living Room", "Child": homeseer.HSLookup{"Name": "Color Control Blue Channel"}}, 0)
		controller.SetChildDevicesValue(homeseer.HSLookup{"LocationOne": "Living Room", "Child": homeseer.HSLookup{"Name": "Color Control Green Channel"}}, 0)
	*/
}
