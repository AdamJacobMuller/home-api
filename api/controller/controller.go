package apicontroller

import (
	"github.com/AdamJacobMuller/home-api/api/models"
	"github.com/AdamJacobMuller/home-api/providers/apc"
	"github.com/AdamJacobMuller/home-api/providers/homeseer"

	log "github.com/Sirupsen/logrus"

	"encoding/json"
	"io/ioutil"
	"os"
)

type ControllerProvider interface {
	SetDevicesValue(apimodels.Match, float64) bool
	SetChildDevicesValue(apimodels.Match, float64) bool

	InvokeDevicesAction(apimodels.Match, string) bool
	InvokeChildDevicesAction(apimodels.Match, string) bool
	TypeString() string
	IDString() string
	Create(json.RawMessage) bool
}

type APIController struct {
	providers []ControllerProvider
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
	for _, provider := range c.providers {
		if c.MatchProvider(match, provider) {
			provider.InvokeDevicesAction(match, action)
		}
	}
	return true
}

type Configuration struct {
	Providers []json.RawMessage `json:"Providers"`
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
		log.WithFields(log.Fields{"filename": filename, "error": err}).Error("unable to unmarshal file")
		return false
	}

	for _, provider = range config.Providers {
		c.CreateProvider(provider)
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
