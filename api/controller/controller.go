package apicontroller

import (
	"github.com/AdamJacobMuller/home-api/api/models"
	"github.com/AdamJacobMuller/home-api/providers/apc"
	"github.com/AdamJacobMuller/home-api/providers/example"
	"github.com/AdamJacobMuller/home-api/providers/homeseer"
	"github.com/AdamJacobMuller/home-api/providers/redeye"

	log "github.com/sirupsen/logrus"

	"encoding/json"
	"io/ioutil"
	"os"
)

type ControllerProvider interface {
	SetDevicesValue(apimodels.Match, float64) bool
	SetChildDevicesValue(apimodels.Match, float64) bool

	InvokeDevicesAction(apimodels.Match, string) bool
	InvokeChildDevicesAction(apimodels.Match, string) bool

	GetDevices(apimodels.Match) (apimodels.Devices, bool)
	GetDevice(apimodels.Match) (apimodels.Device, bool)
	GetChildDevice(apimodels.Match) (apimodels.Device, bool)
	GetChildDevices(apimodels.Match) (apimodels.Devices, bool)

	TypeString() string
	IDString() string
	Create(json.RawMessage) bool
}

type APIController struct {
	providers []ControllerProvider
	triggers  Triggers
}

func (c *APIController) GetDevices(match apimodels.Match) []apimodels.Device {
	var devices []apimodels.Device
	for _, provider := range c.providers {
		provider_devices, ok := provider.GetDevices(match)
		if !ok {
			log.WithFields(log.Fields{"ProviderType": provider.TypeString(), "ProviderID": provider.IDString()}).Error("unable to list devices")
			continue
		}
		for _, device := range provider_devices.List() {
			devices = append(devices, device)
		}
	}
	return devices
}
func (c *APIController) Locations() map[string][]string {
	fs := make(map[string][]string)

	var locationone string
	var locationtwo string
	var doappend bool

	for _, provider := range c.providers {
		devices, ok := provider.GetDevices(apimodels.Match{})
		if !ok {
			log.WithFields(log.Fields{"ProviderType": provider.TypeString(), "ProviderID": provider.IDString()}).Error("unable to list devices")
			continue
		}
		for _, device := range devices.List() {
			locationone = device.GetLocationOne()
			locationtwo = device.GetLocationTwo()
			_, ok := fs[locationtwo]
			if !ok {
				fs[locationtwo] = make([]string, 0)
			}
			doappend = true
			for _, v := range fs[locationtwo] {
				if v == locationone {
					doappend = false
					continue
				}
			}
			if doappend {
				fs[locationtwo] = append(fs[locationtwo], locationone)
			}
		}
	}
	return fs
}
func (c *APIController) MatchProvider(match apimodels.Match, provider ControllerProvider) bool {
	var sVal string
	var ok bool
	sVal, ok = match["ProviderID"].(string)
	if ok {
		if sVal != provider.IDString() {
			return false
		}
	}
	sVal, ok = match["ProviderType"].(string)
	if ok {
		if sVal != provider.TypeString() {
			return false
		}
	}
	return true
}

// SetDevicesValue
// SetChildDevicesValue
// InvokeDevicesAction
// InvokeChildDevicesAction

func (c *APIController) ControlRequest(req apimodels.ControlRequest) bool {
	switch req.Type {
	case "SetDevicesValue":
		return c.SetDevicesValue(req.Match, req.Value)
	case "SetChildDevicesValue":
		return c.SetChildDevicesValue(req.Match, req.Value)
	case "InvokeDevicesAction":
		return c.InvokeDevicesAction(req.Match, req.Action)
	case "InvokeChildDevicesAction":
		return c.InvokeChildDevicesAction(req.Match, req.Action)
	default:
		log.WithFields(log.Fields{"type": req.Type}).Error("invalid ControlRequest type")
		return false
	}
}
func (c *APIController) SetDevicesValue(match apimodels.Match, value float64) bool {
	for _, provider := range c.providers {
		if c.MatchProvider(match, provider) {
			provider.SetDevicesValue(match, value)
		}
	}
	return true
}

func (c *APIController) SetChildDevicesValue(match apimodels.Match, value float64) bool {
	for _, provider := range c.providers {
		if c.MatchProvider(match, provider) {
			provider.SetChildDevicesValue(match, value)
		}
	}
	return true
}

func (c *APIController) InvokeChildDevicesAction(match apimodels.Match, action string) bool {
	for _, provider := range c.providers {
		if c.MatchProvider(match, provider) {
			provider.InvokeChildDevicesAction(match, action)
		}
	}
	return true
}
func (c *APIController) InvokeDevicesAction(match apimodels.Match, action string) bool {
	c.triggers.InvokeDevicesAction(c, "before", match, action)
	for _, provider := range c.providers {
		if c.MatchProvider(match, provider) {
			provider.InvokeDevicesAction(match, action)
		}
	}
	return true
}

type Configuration struct {
	Providers []json.RawMessage `json:"Providers"`
	Triggers  []Trigger         `json:"Triggers"`
}
type ConfigurationProviderType struct {
	ProviderType string `json:"ProviderType"`
}

func (c *APIController) LoadAndCreateProviders(filename string) bool {
	var config *Configuration
	var provider json.RawMessage
	var err error

	config = &Configuration{}

	fh, err := os.Open(filename)
	if err != nil {
		log.WithFields(log.Fields{"filename": filename, "error": err}).Error("unable to open file")
		return false
	}

	data, err := ioutil.ReadAll(fh)
	if err != nil {
		log.WithFields(log.Fields{"filename": filename, "error": err}).Error("unable to read file")
		return false
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		log.WithFields(log.Fields{"filename": filename, "error": err}).Fatal("unable to unmarshal file")
		return false
	}

	c.triggers = Triggers{Triggers: config.Triggers}

	for _, provider = range config.Providers {
		go c.CreateProvider(provider)
	}

	return true
}
func (c *APIController) CreateProvider(providerRaw json.RawMessage) bool {
	var pt *ConfigurationProviderType
	var provider ControllerProvider
	pt = &ConfigurationProviderType{}
	json.Unmarshal(providerRaw, pt)

	log.WithFields(log.Fields{"ProviderType": pt.ProviderType}).Info("creating provider")

	switch pt.ProviderType {
	case "HomeSeer":
		provider = &homeseer.HSController{}
	case "APC PDU":
		provider = &apc.PDU{}
	case "RedEye":
		provider = &redeye.RedEye{}
	case "Example Provider":
		provider = &example.ExampleProvider{}
	default:
		return false
	}
	if provider.Create(providerRaw) {
		log.WithFields(log.Fields{"ProviderType": provider.TypeString(), "ProviderID": provider.IDString()}).Info("create completed")
		c.AddProvider(provider)
		return true
	} else {
		log.WithFields(log.Fields{"ProviderType": provider.TypeString(), "ProviderID": provider.IDString()}).Error("create failed")
		return false
	}
}
func (c *APIController) AddProvider(provider ControllerProvider) {
	c.providers = append(c.providers, provider)
}

func NewAPIController() *APIController {
	c := &APIController{}
	return c
}
