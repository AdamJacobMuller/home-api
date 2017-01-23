package example

import (
	"encoding/json"
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
func (p *ExampleProvider) GetDevices(apimodels.Match) (apimodels.Devices, bool) {
	return &ExampleList{}, false
}
func (p *ExampleProvider) GetDevice(apimodels.Match) (apimodels.Device, bool) {
	return &ExampleDevice{}, false
}
func (p *ExampleProvider) GetChildDevice(apimodels.Match) (apimodels.Device, bool) {
	return &ExampleDevice{}, false
}
func (p *ExampleProvider) GetChildDevices(apimodels.Match) (apimodels.Devices, bool) {
	return &ExampleList{}, false
}
func (p *ExampleProvider) Create(json.RawMessage) bool {
	return false
}
func (p *ExampleProvider) TypeString() string {
	return "ExampleProviderType"
}
func (p *ExampleProvider) IDString() string {
	return "ExampleProviderInstance"
}

type ExampleDevice struct {
}

func (d *ExampleDevice) ListChildren() []apimodels.Device {
	return []apimodels.Device{}
}
func (d *ExampleDevice) ListActions() []apimodels.Action {
	return []apimodels.Action{}
}
func (d *ExampleDevice) GetTypes() []string {
	return []string{}
}
func (d *ExampleDevice) SetValue(float64) bool {
	return false
}
func (d *ExampleDevice) InvokeAction(string) bool {
	return false
}
func (d *ExampleDevice) GetName() string {
	return "ExampleDeviceName"
}
func (d *ExampleDevice) GetLocationOne() string {
	return "Floor"
}
func (d *ExampleDevice) ProviderIDString() string {
	return "ExampleProviderID"
}
func (d *ExampleDevice) GetLocationTwo() string {
	return "Room"
}
func (d *ExampleDevice) IDString() string {
	return "ExampleDeviceID"
}
func (d *ExampleDevice) Matches(apimodels.Match) bool {
	return false
}

type ExampleList struct {
}

func (l *ExampleList) SetValue(float64) bool {
	return false
}
func (l *ExampleList) List() []apimodels.Device {
	return []apimodels.Device{}
}
func (l *ExampleList) InvokeAction(string) bool {
	return false
}
