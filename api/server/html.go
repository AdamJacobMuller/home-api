package apiserver

import (
	"fmt"
	"github.com/AdamJacobMuller/home-api/api/models"
	"github.com/AdamJacobMuller/home-api/api/server/templates"
	"github.com/gocraft/web"
	log "github.com/sirupsen/logrus"
	"sort"
)

type Context struct {
	Layout    *templates.BasePage
	apiserver *APIServer
}

func (c *Context) DrawLayout(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	c.Layout = &templates.BasePage{}
	next(rw, r)
	templates.WritePageTemplate(rw, c.Layout)
}

func (c *Context) AddControllerNavigation(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	c.Layout.MiniNavBar = true
	locations := c.apiserver.Controller.Locations()
	var locationKeys []string
	for ln, _ := range locations {
		locationKeys = append(locationKeys, ln)
	}
	sort.Strings(locationKeys)
	for _, ln := range locationKeys {
		lt := locations[ln]
		sort.Strings(lt)
		sec := c.Layout.NavSection(ln)
		for _, ltn := range lt {
			sec.NavLink(ltn).SetHref(fmt.Sprintf("/room/%s/%s", ln, ltn))
		}
	}
	next(rw, r)
}

func (c *Context) Room(rw web.ResponseWriter, req *web.Request) {
	LocationOne := req.PathParams["LocationOne"]
	LocationTwo := req.PathParams["LocationTwo"]
	rowtemplate := &templates.Rows{}
	c.Layout.Body = rowtemplate
	c.Layout.Title = fmt.Sprintf("Home / %s / %s", LocationOne, LocationTwo)

	boxrow := &templates.GenericRow{}
	rowtemplate.AddRow(boxrow)

	match := apimodels.Match{"LocationOne": LocationOne, "LocationTwo": LocationTwo}
	for _, device := range c.apiserver.Controller.GetDevices(match) {
		if device.IsHidden() {
			continue
		}
		box, ok := DeviceToBox(device)
		if ok {
			boxrow.AddBox(box)
		}
	}
}

func DeviceToBox(device apimodels.Device) (templates.Box, bool) {
	for i, devicetype := range device.ListTypes() {
		switch devicetype.GetName() {
		case "Z-Wave Entry Control Root Device":
			return DoorLockBox(device)
		case "Z-Wave Switch Multilevel Root Device":
			return ColorChangeBulbBox(device, false)
		case "Sonos Player Master Control":
			return SonosBox(device)
		case "Z-Wave Switch Binary":
			if device.HasChildDevice(apimodels.Match{"TypeString": "Z-Wave Switch Multilevel Root Device"}) {
				return ColorChangeBulbBox(device, true)
			} else {
				return BinarySwitchBox(device)
			}
		case "Z-Wave Switch Binary Root Device":
			return BinarySwitchBox(device)
		case "Z-Wave Switch Multilevel":
			return DimmableSwitchBox(device)
		case "TiVo":
			return TivoBox(device)
		case "RedEye":
			return RedEyeBox(device)
		case "Thermostat":
			return Thermostat(device)
		default:
			log.WithFields(log.Fields{"i": i, "type": devicetype, "name": devicetype.GetName(), "id": device.IDString()}).Error("unable to locate matching device type")
		}
	}
	return GenericBox(device)
}
func ColorChangeBulbBox(device apimodels.Device, PowerControlRoot bool) (templates.Box, bool) {
	return templates.ColorChangeBulbBox{Title: device.GetName(), DeviceID: device.IDString(), ProviderID: device.ProviderIDString(), PowerControlRoot: PowerControlRoot}, true
}
func DimmableSwitchBox(device apimodels.Device) (templates.Box, bool) {
	return templates.DimmableSwitchBox{Title: device.GetName(), DeviceID: device.IDString(), ProviderID: device.ProviderIDString()}, true
}
func BinarySwitchBox(device apimodels.Device) (templates.Box, bool) {
	return templates.BinarySwitchBox{Title: device.GetName(), DeviceID: device.IDString(), ProviderID: device.ProviderIDString()}, true
}
func SonosBox(device apimodels.Device) (templates.Box, bool) {
	return templates.SonosBox{Title: device.GetName(), DeviceID: device.IDString(), ProviderID: device.ProviderIDString()}, true
}
func TivoBox(device apimodels.Device) (templates.Box, bool) {
	return templates.TivoBox{Title: device.GetName(), DeviceID: device.IDString(), ProviderID: device.ProviderIDString()}, true
}
func Thermostat(device apimodels.Device) (templates.Box, bool) {
	box := templates.Thermostat{Title: device.GetName(), DeviceID: device.IDString(), ProviderID: device.ProviderIDString()}
	for _, action := range device.ListActions() {
		boxAction := &templates.Action{Title: action.GetName()}
		box.AddAction(boxAction)
	}
	return box, true
}
func RedEyeBox(device apimodels.Device) (templates.Box, bool) {
	box := templates.RedEyeBox{Title: device.GetName(), DeviceID: device.IDString(), ProviderID: device.ProviderIDString()}
	for _, action := range device.ListActions() {
		boxAction := &templates.Action{Title: action.GetName()}
		box.AddAction(boxAction)
	}
	return box, true
}
func DoorLockBox(device apimodels.Device) (templates.Box, bool) {
	return templates.DoorLockBox{Title: device.GetName(), DeviceID: device.IDString(), ProviderID: device.ProviderIDString()}, true
}
func GenericBox(device apimodels.Device) (templates.Box, bool) {
	var includeDevice bool
	var includeChildDevice bool
	includeDevice = false
	box := &templates.DeviceBox{}
	box.Title = device.GetName()
	for _, action := range device.ListActions() {
		boxAction := &templates.Action{Title: action.GetName()}
		box.AddAction(boxAction)
		includeDevice = true
	}
	for _, child := range device.ListChildren() {
		includeChildDevice = false
		childBox := &templates.ChildDeviceBox{}
		childBox.Title = child.GetName()
		for _, action := range child.ListActions() {
			boxAction := &templates.Action{Title: action.GetName()}
			childBox.AddAction(boxAction)
			includeDevice = true
			includeChildDevice = true
		}
		if includeChildDevice {
			box.AddChild(childBox)
		}
	}
	if includeDevice {
		return box, true
	} else {
		return nil, false
	}
}

func (c *Context) Index(rw web.ResponseWriter, req *web.Request) {
	fmt.Printf("we are in the index\n")
}
func (c *Context) NotFound(rw web.ResponseWriter, req *web.Request) {
	rw.WriteHeader(404)
	c.Layout.Body = &templates.NotFound{}
}
