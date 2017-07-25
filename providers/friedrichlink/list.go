package friedrichlink

import (
	"github.com/AdamJacobMuller/home-api/api/models"
)

type FriedrichLinkList struct {
	devices []*FriedrichLinkDevice
}

func (l *FriedrichLinkList) SetValue(float64) bool {
	return false
}
func (l *FriedrichLinkList) List() []apimodels.Device {
	var v []apimodels.Device
	for _, i := range l.devices {
		v = append(v, i)
	}
	return v
}
func (l *FriedrichLinkList) InvokeAction(string) bool {
	return false
}
func (l *FriedrichLinkList) Add(device *FriedrichLinkDevice) {
	l.devices = append(l.devices, device)
}
