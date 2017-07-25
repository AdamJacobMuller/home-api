package friedrichlink

import (
	"github.com/AdamJacobMuller/home-api/api/models"
)

type FriedrichLinkDevice struct {
	Name      string
	Serial    string
	Ambient   string
	Model     string
	Power     string
	CoolPoint string
	fl        *FriedrichLink
}

func (d *FriedrichLinkDevice) HasChildDevice(find apimodels.Match) bool {
	return false
}
func (d *FriedrichLinkDevice) GetChildDevice(find apimodels.Match) apimodels.Device {
	return &FriedrichLinkDevice{}
}
func (d *FriedrichLinkDevice) IsHidden() bool {
	return false
}
func (d *FriedrichLinkDevice) ListChildren() []apimodels.Device {
	return []apimodels.Device{}
}
func (d *FriedrichLinkDevice) ListActions() []apimodels.Action {
	return []apimodels.Action{}
}
func (d *FriedrichLinkDevice) ListTypes() []apimodels.Type {
	return []apimodels.Type{Thermostat{}}
}
func (d *FriedrichLinkDevice) SetValue(float64) bool {
	return false
}
func (d *FriedrichLinkDevice) InvokeAction(string) bool {
	return false
}
func (d *FriedrichLinkDevice) GetName() string {
	return d.Name
}
func (d *FriedrichLinkDevice) GetLocationOne() string {
	return "Floor"
}
func (d *FriedrichLinkDevice) ProviderIDString() string {
	return d.fl.UserID
}
func (d *FriedrichLinkDevice) GetLocationTwo() string {
	return "Room"
}
func (d *FriedrichLinkDevice) IDString() string {
	return d.Serial
}
func (d *FriedrichLinkDevice) Matches(apimodels.Match) bool {
	return true
}
