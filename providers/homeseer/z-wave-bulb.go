package homeseer

import (
	"github.com/AdamJacobMuller/home-api/api/models"
)

func ZWSML_Red(device *HSDevice, action *HSAction) bool {
	device.GetChildDevice(apimodels.Match{"Name": "Color Control Red Channel"}).SetValue(255)
	device.GetChildDevice(apimodels.Match{"Name": "Color Control Green Channel"}).SetValue(0)
	device.GetChildDevice(apimodels.Match{"Name": "Color Control Blue Channel"}).SetValue(0)
	device.GetChildDevice(apimodels.Match{"Name": "Color Control Warm_White Channel"}).SetValue(0)
	device.GetChildDevice(apimodels.Match{"Name": "Color Control Cold_White Channel"}).SetValue(0)
	return true
}

func ZWSML_Green(device *HSDevice, action *HSAction) bool {
	device.GetChildDevice(apimodels.Match{"Name": "Color Control Red Channel"}).SetValue(0)
	device.GetChildDevice(apimodels.Match{"Name": "Color Control Green Channel"}).SetValue(255)
	device.GetChildDevice(apimodels.Match{"Name": "Color Control Blue Channel"}).SetValue(0)
	device.GetChildDevice(apimodels.Match{"Name": "Color Control Warm_White Channel"}).SetValue(0)
	device.GetChildDevice(apimodels.Match{"Name": "Color Control Cold_White Channel"}).SetValue(0)
	return true
}

func ZWSML_Blue(device *HSDevice, action *HSAction) bool {
	device.GetChildDevice(apimodels.Match{"Name": "Color Control Red Channel"}).SetValue(0)
	device.GetChildDevice(apimodels.Match{"Name": "Color Control Green Channel"}).SetValue(0)
	device.GetChildDevice(apimodels.Match{"Name": "Color Control Blue Channel"}).SetValue(255)
	device.GetChildDevice(apimodels.Match{"Name": "Color Control Warm_White Channel"}).SetValue(0)
	device.GetChildDevice(apimodels.Match{"Name": "Color Control Cold_White Channel"}).SetValue(0)
	return true
}

func ZWSML_Warm_White(device *HSDevice, action *HSAction) bool {
	device.GetChildDevice(apimodels.Match{"Name": "Color Control Red Channel"}).SetValue(0)
	device.GetChildDevice(apimodels.Match{"Name": "Color Control Green Channel"}).SetValue(0)
	device.GetChildDevice(apimodels.Match{"Name": "Color Control Blue Channel"}).SetValue(0)
	device.GetChildDevice(apimodels.Match{"Name": "Color Control Warm_White Channel"}).SetValue(255)
	device.GetChildDevice(apimodels.Match{"Name": "Color Control Cold_White Channel"}).SetValue(0)
	return true
}

func ZWSML_Cold_White(device *HSDevice, action *HSAction) bool {
	device.GetChildDevice(apimodels.Match{"Name": "Color Control Red Channel"}).SetValue(0)
	device.GetChildDevice(apimodels.Match{"Name": "Color Control Green Channel"}).SetValue(0)
	device.GetChildDevice(apimodels.Match{"Name": "Color Control Blue Channel"}).SetValue(0)
	device.GetChildDevice(apimodels.Match{"Name": "Color Control Warm_White Channel"}).SetValue(0)
	device.GetChildDevice(apimodels.Match{"Name": "Color Control Cold_White Channel"}).SetValue(255)
	return true
}
