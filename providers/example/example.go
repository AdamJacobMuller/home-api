package example

import (
	"github.com/AdamJacobMuller/home-api/api/models"
)

type ExampleProvider struct {
}

func (p *ExampleProvider) SetDevicesValue(apimodels.Match, float64) bool {
	return false
}
func (p *ExampleProvider) SetChildDevicesValue(apimodels.Match, float64) bool {
	return false
}
func (p *ExampleProvider) InvokeDevicesAction(apimodels.Match, string) bool {
	return false
}
func (p *ExampleProvider) InvokeChildDevicesAction(apimodels.Match, string) bool {
	return false
}
func (p *ExampleProvider) TypeString() string {
	return "ExampleProviderType"
}
func (p *ExampleProvider) IDString() string {
	return "ExampleProviderInstance"
}
