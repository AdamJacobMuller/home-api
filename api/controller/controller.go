package apicontroller

import (
	"github.com/AdamJacobMuller/home-api/api/models"
)

type ControllerProvider interface {
	SetDevicesValue(apimodels.Match, float64) bool
	SetChildDevicesValue(apimodels.Match, float64) bool

	InvokeDevicesAction(apimodels.Match, string) bool
	InvokeChildDevicesAction(apimodels.Match, string) bool
	IDString() string
}

type APIController struct {
	providers []ControllerProvider
}

func (c *APIController) SetDevicesValue(match apimodels.Match, value float64) bool {
	for _, provider := range c.providers {
		provider.SetDevicesValue(match, value)
	}
	return true
}

func (c *APIController) SetChildDevicesValue(match apimodels.Match, value float64) bool {
	for _, provider := range c.providers {
		provider.SetChildDevicesValue(match, value)
	}
	return true
}

func (c *APIController) InvokeChildDevicesAction(match apimodels.Match, action string) bool {
	for _, provider := range c.providers {
		provider.InvokeChildDevicesAction(match, action)
	}
	return true
}
func (c *APIController) InvokeDevicesAction(match apimodels.Match, action string) bool {
	for _, provider := range c.providers {
		provider.InvokeDevicesAction(match, action)
	}
	return true
}

func (c *APIController) AddProvider(provider ControllerProvider) {
	c.providers = append(c.providers, provider)
}

func NewAPIController() *APIController {
	c := &APIController{}
	return c
}
