package apiserver

import (
	"fmt"

	"github.com/AdamJacobMuller/home-api/api/models"
	"github.com/gocraft/web"
)

func (c *Context) Room(w web.ResponseWriter, req *web.Request) {
	LocationOne := req.PathParams["LocationOne"]
	LocationTwo := req.PathParams["LocationTwo"]

	match := apimodels.Match{"LocationOne": LocationOne, "LocationTwo": LocationTwo}
	for _, device := range c.apiserver.Controller.GetDevices(match) {
		if device.IsHidden() {
			continue
		}
		fmt.Printf("Name: %s\n", device.GetName())
		for _, deviceType := range device.ListTypes() {
			fmt.Printf("    Type: %s\n", deviceType)
		}
		for _, action := range device.ListActions() {
			fmt.Printf("    Acti: %s\n", action.GetName())
		}
		for i, child := range device.ListChildren() {
			fmt.Printf("    [%d]Name: %s\n", i, child.GetName())
			for _, deviceType := range child.ListTypes() {
				fmt.Printf("    [%d]    Type: %s\n", i, deviceType)
			}
			for _, action := range child.ListActions() {
				fmt.Printf("    [%d]    Acti: %s\n", i, action.GetName())
			}
		}
	}
}
