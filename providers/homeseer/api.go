package homeseer

import (
	"encoding/json"
	"fmt"
	"github.com/AdamJacobMuller/home-api/api/models"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type API struct {
	Base string
}

func (c *API) GetAndUnmarshal(path string, object interface{}) error {
	url := fmt.Sprintf("%s/%s", c.Base, path)
	log.WithFields(log.Fields{"url": url}).Info("GetAndUnmarshal")
	resp, err := http.Get(url)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("GET failed")
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("read Body failed")
		return err
	}
	resp.Body.Close()

	err = json.Unmarshal(body, object)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("json.Unmarshal failed")
		return err
	}
	return nil
}

type HSAction struct {
	API      *API
	DeviceID int
	Label    string
	Value    float64
}
type HSDevice struct {
	ID          int
	Name        string
	API         *API
	LocationOne string
	LocationTwo string
	Value       float64
	Status      string
	TypeString  string
	LastChange  *time.Time
	Children    []*HSDevice
	Parent      *HSDevice
	Actions     []*HSAction
}

type HSController struct {
	API            *API
	UpdateInterval time.Duration
	Devices        []*HSDevice
}

func (c *HSController) IDString() string {
	return c.API.Base
}

func (h *HSDevice) InvokeAction(label string) bool {
	url := fmt.Sprintf("JSON?request=controldevicebylabel&ref=%d&label=%s", h.ID, label)
	json_devices := &JD_HSDevices{}
	err := h.API.GetAndUnmarshal(url, json_devices)
	if err != nil {
		log.WithFields(log.Fields{"id": h.ID, "error": err}).Error("unable to controldevicebyvalue")
		return false
	}
	if len(json_devices.Devices) == 1 {
		// lol?
		h = json_devices.Devices[0].asHS()
		return true
	} else {
		log.WithFields(log.Fields{"id": h.ID, "reply": json_devices}).Error("unable to controldevicebyvalue")
		return false
	}
}

func (h *HSDevice) SetValue(value float64) bool {
	url := fmt.Sprintf("JSON?request=controldevicebyvalue&ref=%d&value=%f", h.ID, value)
	json_devices := &JD_HSDevices{}
	err := h.API.GetAndUnmarshal(url, json_devices)
	if err != nil {
		log.WithFields(log.Fields{"id": h.ID, "error": err}).Error("unable to controldevicebyvalue")
		return false
	}
	if len(json_devices.Devices) == 1 {
		// lol?
		h = json_devices.Devices[0].asHS()
		return true
	} else {
		log.WithFields(log.Fields{"id": h.ID, "reply": json_devices}).Error("unable to controldevicebyvalue")
		return false
	}
}
func (h *HSDevice) IDString() string {
	return fmt.Sprintf("%d", h.ID)
}
func (h *HSDevice) Matches(find apimodels.Match) bool {
	var sVal string
	var lVal apimodels.Match
	var ok bool
	sVal, ok = find["LocationOne"].(string)
	if ok {
		if sVal != h.LocationOne {
			return false
		}
	}
	sVal, ok = find["LocationTwo"].(string)
	if ok {
		if sVal != h.LocationTwo {
			return false
		}
	}
	sVal, ok = find["Name"].(string)
	if ok {
		if sVal != h.Name {
			return false
		}
	}
	sVal, ok = find["TypeString"].(string)
	if ok {
		if sVal != h.TypeString {
			return false
		}
	}
	lVal, ok = find["Child"].(apimodels.Match)
	if ok {
		var matched bool
		matched = false
		for _, child := range h.Children {
			if child.Matches(lVal) {
				matched = true
			}
		}
		if !matched {
			return false
		}
	}
	return true
}

type HSResult struct {
	Devices []*HSDevice
}

func (h *HSResult) Print() {
	for _, device := range h.Devices {
		device.Print()
	}
}
func (h *HSResult) InvokeAction(action string) bool {
	for _, device := range h.Devices {
		device.InvokeAction(action)
	}
	return true
}
func (h *HSResult) SetValue(value float64) bool {
	for _, device := range h.Devices {
		device.SetValue(value)
	}
	return true
}

func (h *HSResult) Add(device *HSDevice) {
	h.Devices = append(h.Devices, device)
}

func (h *HSController) SetChildDevicesValue(find apimodels.Match, value float64) bool {
	devices, ok := h.GetChildDevices(find)
	if ok {
		devices.SetValue(value)
		return true
	}
	return false
}
func (h *HSController) InvokeChildDevicesAction(find apimodels.Match, action string) bool {
	devices, ok := h.GetChildDevices(find)
	if ok {
		devices.InvokeAction(action)
		return true
	}
	return false
}
func (h *HSController) InvokeDevicesAction(find apimodels.Match, action string) bool {
	devices, ok := h.GetDevices(find)
	if ok {
		devices.InvokeAction(action)
		return true
	}
	return false
}
func (h *HSController) SetDevicesValue(find apimodels.Match, value float64) bool {
	devices, ok := h.GetDevices(find)
	if ok {
		devices.SetValue(value)
		return true
	}
	return false
}
func (h *HSController) GetDevices(find apimodels.Match) (apimodels.Devices, bool) {
	result := &HSResult{}
	for _, device := range h.Devices {
		if device.Matches(find) {
			result.Add(device)
		}
	}
	return result, true
}
func (h *HSController) GetDevice(find apimodels.Match) (apimodels.Device, bool) {
	for _, device := range h.Devices {
		if device.Matches(find) {
			return device, true
		}
	}
	return &HSDevice{}, false
}

func (h *HSController) GetChildDevice(find apimodels.Match) (apimodels.Device, bool) {
	var cLookup apimodels.Match
	var ok bool
	cLookup, ok = find["Child"].(map[string]interface{})
	if !ok {
		log.Error("GetChildDevice apimodels.Match requires a Child apimodels.Match element")
		return &HSDevice{}, false
	}
	for _, device := range h.Devices {
		if device.Matches(find) {
			for _, childDevice := range device.Children {
				if childDevice.Matches(cLookup) {
					return childDevice, true
				}
			}
		}
	}
	return &HSDevice{}, false
}
func (h *HSController) GetChildDevices(find apimodels.Match) (apimodels.Devices, bool) {
	result := &HSResult{}
	var cLookup apimodels.Match
	var ok bool
	cLookup, ok = find["Child"].(map[string]interface{})
	if !ok {
		log.Error("GetChildDevice apimodels.Match requires a Child apimodels.Match element")
		return result, false
	}
	for _, device := range h.Devices {
		if device.Matches(find) {
			for _, childDevice := range device.Children {
				if childDevice.Matches(cLookup) {
					result.Add(childDevice)
				}
			}
		}
	}
	return result, true
}

func (h *HSController) Print() {
	for _, device := range h.Devices {
		device.Print()
	}
}

func (h *HSDevice) Print() {
	fmt.Printf("%s [Value = %f]\n", h.Name, h.Value)
	for _, action := range h.Actions {
		fmt.Printf(" - %s [Value = %f]\n", action.Label, action.Value)
	}
	for _, child := range h.Children {
		fmt.Printf("    %s [Value = %f]\n", child.Name, child.Value)
		for _, action := range child.Actions {
			fmt.Printf("     - %s [Value = %f]\n", action.Label, action.Value)
		}
	}
}
func (h *HSDevice) SetControls(actions []*HSAction) {
	h.Actions = actions
}
func (h *HSDevice) AddChild(device *HSDevice) {
	device.Parent = h
	h.Children = append(h.Children, device)
}

func (h *HSController) Load() {
	json_devices := &JD_HSDevices{}
	json_control_devices := &JD_HSControlDevices{}

	h.API.GetAndUnmarshal("JSON?request=getstatus", json_devices)
	h.API.GetAndUnmarshal("JSON?request=getcontrol&ref=all", json_control_devices)
	hold_devices := make(map[int][]*HSDevice)
	devices := make(map[int]*HSDevice)
	all_controls := make(map[int][]*HSAction)

	for _, json_control_device := range json_control_devices.Devices {
		for _, pair := range json_control_device.ControlPairs {
			hsp := pair.asHS()
			hsp.API = h.API
			all_controls[json_control_device.ID] = append(all_controls[json_control_device.ID], hsp)
		}
	}

	// relationship =
	// 2:root device (other devices may be part of this physical device)
	// 3:standalone=this is the only device that represents this physical device
	// 4:child=this device is part of a group of devices that represent this physical device

	for _, device := range json_devices.Devices {
		nd := device.asHS()
		nd.API = h.API
		controls, ok := all_controls[device.ID]
		if ok {
			nd.SetControls(controls)
		}
		switch device.Relationship {
		case 2:
			devices[device.ID] = nd
		case 3:
			devices[device.ID] = nd
		case 4:
			parent_id := device.Associated_devices[0]
			pdevice, ok := devices[parent_id]
			if ok {
				pdevice.AddChild(nd)
			} else {
				hold_devices[parent_id] = append(hold_devices[parent_id], nd)
			}
		}
	}
	for parent_id, device_l := range hold_devices {
		for _, device := range device_l {
			pdevice, ok := devices[parent_id]
			if ok {
				pdevice.AddChild(device)
			} else {
				log.WithFields(log.Fields{"parent_id": parent_id, "id": device.ID}).Fatal("unable to find parent device")
			}
		}
	}
	listDevices := []*HSDevice{}
	for _, device := range devices {
		listDevices = append(listDevices, device)
	}
	h.Devices = listDevices
}
func (h *HSController) BGUpdate() {
	for {
		h.Load()
		time.Sleep(h.UpdateInterval)
	}
}

func NewHomeseerController(api string) *HSController {
	controller := &HSController{
		API:            &API{Base: api},
		UpdateInterval: time.Duration(30 * time.Second),
	}
	controller.Load()
	go controller.BGUpdate()
	return controller
}
