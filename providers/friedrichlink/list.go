package friedrichlink

import (
	"github.com/AdamJacobMuller/home-api/api/models"
)

type FriedrichLinkList struct {
}

func (l *FriedrichLinkList) SetValue(float64) bool {
	return false
}
func (l *FriedrichLinkList) List() []apimodels.Device {
	return []apimodels.Device{}
}
func (l *FriedrichLinkList) InvokeAction(string) bool {
	return false
}
