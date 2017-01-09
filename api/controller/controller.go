package apicontroller

import (
	"github.com/AdamJacobMuller/home-api/api/models"
)

type ControllerProvider interface {
	SetDevicesValue(apimodels.Match, float64) bool
	SetChildDevicesValue(apimodels.Match, float64) bool
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

func (c *APIController) AddProvider(provider ControllerProvider) {
	c.providers = append(c.providers, provider)
}

func NewAPIController() *APIController {
	c := &APIController{}
	return c
}
