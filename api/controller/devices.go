package apicontroller

import (
	"github.com/AdamJacobMuller/home-api/api/models"

	log "github.com/Sirupsen/logrus"
)

type JSON_Devices struct {
	Devices []JSON_Device `json:"Devices"`
}

func (l *JSON_Devices) Add(jd JSON_Device) {
	l.Devices = append(l.Devices, jd)
}

type JSON_Action struct {
	Name string `json:"Name"`
}

type JSON_Type struct {
	Name string `json:"Name"`
}

type JSON_Device struct {
	Name        string        `json:"Name"`
	LocationOne string        `json:"LocationOne"`
	LocationTwo string        `json:"LocationTwo"`
	Types       []JSON_Type   `json:"Types"`
	DeviceID    string        `json:"DeviceID"`
	ProviderID  string        `json:"ProviderID"`
	Children    []JSON_Device `json:"Children,omitempty"`
	Actions     []JSON_Action `json:"Actions,omitempty"`
	Hidden      bool          `json:"Hidden"`
}

func (d *JSON_Device) AddChild(device JSON_Device) {
	d.Children = append(d.Children, device)
}

func (d *JSON_Device) AddType(jtype apimodels.Type) {
	d.Types = append(d.Types, JSON_Type{Name: jtype.GetName()})
}

func (d *JSON_Device) AddAction(action apimodels.Action) {
	d.Actions = append(d.Actions, JSON_Action{Name: action.GetName()})
}

func DeviceTOJSON(device apimodels.Device) (json_device JSON_Device) {
	json_device = JSON_Device{
		Name:        device.GetName(),
		DeviceID:    device.IDString(),
		ProviderID:  device.ProviderIDString(),
		LocationOne: device.GetLocationOne(),
		LocationTwo: device.GetLocationTwo(),
		Hidden:      device.IsHidden(),
	}

	for _, typestring := range device.ListTypes() {
		json_device.AddType(typestring)
	}

	for _, child := range device.ListChildren() {
		json_device.AddChild(DeviceTOJSON(child))
	}

	for _, action := range device.ListActions() {
		json_device.AddAction(action)
	}

	return
}

func (c *APIController) ListDevices() JSON_Devices {
	j_devices := JSON_Devices{}
	for _, provider := range c.providers {
		devices, ok := provider.GetDevices(apimodels.Match{})
		if !ok {
			log.WithFields(log.Fields{"ProviderType": provider.TypeString(), "ProviderID": provider.IDString()}).Error("unable to list devices")
		}
		for _, device := range devices.List() {
			j_devices.Add(DeviceTOJSON(device))
		}
	}
	return j_devices
}
