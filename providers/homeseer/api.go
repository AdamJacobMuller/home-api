package homeseer

import (
	"encoding/json"
	"fmt"
	"github.com/AdamJacobMuller/home-api/api/models"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"regexp"
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

type HSType struct {
	Name string
}

func (t *HSType) GetName() string {
	return t.Name
}

type HSActionFunction func(*HSDevice, *HSAction) bool

type HSAction struct {
	API      *API
	DeviceID int
	Label    string
	Value    float64
	Function HSActionFunction
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
	DeviceType  *JD_HSDeviceType
	Actions     []*HSAction
	Hidden      bool
}

func (d *HSDevice) IsHidden() bool {
	return d.Hidden
}

func (d *HSDevice) ListTypes() []apimodels.Type {
	var r []apimodels.Type

	if d.TypeString != "" {
		r = append(r, &HSType{Name: d.TypeString})
	}
	if d.DeviceType.Device_SubType_Description != "" {
		r = append(r, &HSType{Name: d.DeviceType.Device_SubType_Description})
	}
	if d.DeviceType.Device_Type_Description != "" {
		r = append(r, &HSType{Name: d.DeviceType.Device_Type_Description})
	}
	if d.DeviceType.Device_API_Description != "" {
		r = append(r, &HSType{Name: d.DeviceType.Device_API_Description})
	}
	return r
}
func (a *HSAction) GetName() string {
	return a.Label
}

func (h *HSDevice) ListActions() []apimodels.Action {
	l := make([]apimodels.Action, 0)
	for _, a := range h.Actions {
		l = append(l, a)
	}
	return l
}

type HSController struct {
	API            *API
	UpdateInterval time.Duration
	Devices        []*HSDevice
	ChildMapping   []ChildMapping
	HideDevices    []apimodels.Match
}

func (c *HSController) IDString() string {
	return c.API.Base
}
func (c *HSController) TypeString() string {
	return "HomeSeer"
}

func (h *HSDevice) GetLocationOne() string {
	return h.LocationOne
}
func (h *HSDevice) GetLocationTwo() string {
	return h.LocationTwo
}
func (h *HSDevice) ListChildren() []apimodels.Device {
	r := make([]apimodels.Device, 0)
	for _, c := range h.Children {
		r = append(r, c)
	}
	return r
}

func (h *HSDevice) InvokeAction(label string) bool {
	for _, action := range h.Actions {
		log.WithFields(log.Fields{"Invoke/Label": label, "Action.Label": action.Label, "Action.Function": action.Function}).Info("searching for function")
		if action.Label == label && action.Function != nil {
			return action.Function(h, action)
		}
	}
	return h.RealInvokeAction(label)
}
func (h *HSDevice) RealInvokeAction(label string) bool {
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
func (h *HSDevice) GetName() string {
	return h.Name
}
func (h *HSDevice) ProviderIDString() string {
	return h.API.Base
}
func (h *HSDevice) IDString() string {
	return fmt.Sprintf("%d", h.ID)
}
func (h *HSDevice) Matches(find apimodels.Match) bool {
	var rVal *regexp.Regexp
	var sVal string
	var lVal apimodels.Match
	var ok bool
	var err error
	sVal, ok = find["RegexpName"].(string)
	if ok {
		rVal, err = regexp.Compile(sVal)
		if err != nil {
			log.WithFields(log.Fields{"regexp": sVal, "error": err}).Error("unable to compile regexp")
			return false
		}
		ok = rVal.MatchString(h.Name)
		if !ok {
			return false
		}
	}
	sVal, ok = find["DeviceID"].(string)
	if ok {
		if sVal != h.IDString() {
			return false
		}
	}
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

func (h *HSResult) List() []apimodels.Device {
	r := make([]apimodels.Device, 0)
	for _, c := range h.Devices {
		r = append(r, c)
	}
	return r
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

func getDevices(devices []*HSDevice, find apimodels.Match) (*HSResult, bool) {
	result := &HSResult{}
	for _, device := range devices {
		if device.Matches(find) {
			result.Add(device)
		}
	}
	return result, true
}

func getDevice(devices []*HSDevice, find apimodels.Match) (*HSDevice, bool) {
	for _, device := range devices {
		if device.Matches(find) {
			return device, true
		}
	}
	return &HSDevice{}, false
}

func (h *HSController) GetDevice(find apimodels.Match) (apimodels.Device, bool) {
	for _, device := range h.Devices {
		if device.Matches(find) {
			return device, true
		}
	}
	return &HSDevice{}, false
}

func (d *HSDevice) HasChildDevice(find apimodels.Match) bool {
	for _, device := range d.Children {
		if device.Matches(find) {
			return true
		}
	}
	return false
}

func (d *HSDevice) GetChildDevices(find apimodels.Match) (apimodels.Devices, bool) {
	result := &HSResult{}
	for _, device := range d.Children {
		if device.Matches(find) {
			result.Add(device)
		}
	}
	return result, true
}
func (d *HSDevice) GetChildDevice(find apimodels.Match) apimodels.Device {
	for _, device := range d.Children {
		if device.Matches(find) {
			return device
		}
	}
	return &HSDevice{}
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
func (h *HSDevice) AddActionFunction(label string, function HSActionFunction) {
	h.Actions = append(h.Actions, &HSAction{Label: label, Function: function})
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

	for _, find := range h.HideDevices {
		devices, _ := getDevices(listDevices, find)
		for _, device := range devices.Devices {
			log.WithFields(log.Fields{"device": device}).Info("hiding device")
			device.Hidden = true
		}
	}
	for _, mapping := range h.ChildMapping {
		parents, _ := getDevices(listDevices, mapping.Parent)
		for _, parent := range parents.Devices {
			for _, childMatch := range mapping.Children {
				children, _ := getDevices(listDevices, childMatch)
				for _, child := range children.Devices {
					if mapping.HideChildren {
						child.Hidden = true
					}
					//log.WithFields(log.Fields{"parent": parent, "child": child}).Info("adding Child to Parent")
					parent.AddChild(child)
				}
			}
		}
	}

	for _, device := range listDevices {
		if device.HasChildDevice(apimodels.Match{"TypeString": "Z-Wave Switch Multilevel Root Device"}) {
			device.AddActionFunction("On", Child_ZWSML_On)
			device.AddActionFunction("Off", Child_ZWSML_Off)
			device.AddActionFunction("Red", Child_ZWSML_Red)
			device.AddActionFunction("Green", Child_ZWSML_Green)
			device.AddActionFunction("Blue", Child_ZWSML_Blue)
			device.AddActionFunction("Warm White", Child_ZWSML_Warm_White)
			device.AddActionFunction("Cold White", Child_ZWSML_Cold_White)
			device.AddActionFunction("White", Child_ZWSML_Warm_White)
		}

		for _, devicetype := range device.ListTypes() {
			switch devicetype.GetName() {
			case "Z-Wave Switch Multilevel Root Device":
				device.AddActionFunction("Red", ZWSML_Red)
				device.AddActionFunction("Green", ZWSML_Green)
				device.AddActionFunction("Blue", ZWSML_Blue)
				device.AddActionFunction("Warm White", ZWSML_Warm_White)
				device.AddActionFunction("Cold White", ZWSML_Cold_White)
				device.AddActionFunction("White", ZWSML_Warm_White)
			}
		}
	}

	h.Devices = listDevices
}
func (h *HSController) BGUpdate() {
	for {
		h.Load()
		time.Sleep(h.UpdateInterval)
	}
}

type HSConfiguration struct {
	API            string            `json:"API"`
	UpdateInterval int               `json:"UpdateInterval"`
	ChildMapping   []ChildMapping    `json:"ChildMapping"`
	HideDevices    []apimodels.Match `json:"HideDevices"`
}

type ChildMapping struct {
	Parent       apimodels.Match   `json:"Parent"`
	Children     []apimodels.Match `json:"Children"`
	HideChildren bool              `json:"HideChildren"`
}

func (h *HSController) Create(configurationRaw json.RawMessage) bool {
	var configuration *HSConfiguration

	configuration = &HSConfiguration{}

	json.Unmarshal(configurationRaw, configuration)

	if configuration.API == "" {
		log.WithFields(log.Fields{}).Error("API is a required configuration option")
		return false
	} else {
		h.API = &API{
			Base: configuration.API,
		}
	}
	if configuration.UpdateInterval == 0 {
		h.UpdateInterval = time.Duration(30 * time.Second)
	} else {
		h.UpdateInterval = time.Duration(time.Duration(configuration.UpdateInterval) * time.Second)
	}
	h.ChildMapping = configuration.ChildMapping
	h.HideDevices = configuration.HideDevices
	h.Load()
	go h.BGUpdate()
	return true
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
