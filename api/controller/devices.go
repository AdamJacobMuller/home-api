package apicontroller

import (
	"github.com/AdamJacobMuller/home-api/api/models"

	log "github.com/Sirupsen/logrus"
)

type JSON_Devices struct {
	Devices []JSON_Device `json:"devices"`
}

type JSON_Device struct {
	Name string `json:"name"`
	ID   string `json:"string"`
}

func (c *APIController) ListDevices() JSON_Devices {
	j_devices := JSON_Devices{}
	for _, provider := range c.providers {
		devices, ok := provider.GetDevices(apimodels.Match{})
		if !ok {
			log.WithFields(log.Fields{"ProviderType": provider.TypeString(), "ProviderID": provider.IDString()}).Error("unable to list devices")
		}
		for _, device := range devices.List() {
			j_devices.Devices = append(j_devices.Devices, JSON_Device{
				Name: device.GetName(),
				ID:   device.IDString(),
			})
		}
	}
	return j_devices
}
