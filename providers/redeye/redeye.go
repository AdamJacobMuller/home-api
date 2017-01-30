package redeye

import (
	"encoding/json"
	"github.com/AdamJacobMuller/home-api/api/models"
)

type RedEye struct {
	API   *API
	Rooms []*RedEyeRoom
}

func (p *RedEye) SetDevicesValue(apimodels.Match, float64) bool {
	return false
}
func (p *RedEye) SetChildDevicesValue(apimodels.Match, float64) bool {
	return false
}
func (p *RedEye) InvokeDevicesAction(find apimodels.Match, action string) bool {
	for _, room := range p.Rooms {
		if room.Matches(find) {
			room.InvokeAction(action)
		}
	}
	return false
}
func (p *RedEye) InvokeChildDevicesAction(apimodels.Match, string) bool {
	return false
}
func (p *RedEye) GetDevices(find apimodels.Match) (apimodels.Devices, bool) {
	list := &DeviceList{}
	for _, room := range p.Rooms {
		if room.Matches(find) {
			list.Add(room)
		}
	}
	return list, true
}
func (p *RedEye) GetDevice(apimodels.Match) (apimodels.Device, bool) {
	return &Device{}, false
}
func (p *RedEye) GetChildDevice(apimodels.Match) (apimodels.Device, bool) {
	return &Device{}, false
}
func (p *RedEye) GetChildDevices(apimodels.Match) (apimodels.Devices, bool) {
	return &DeviceList{}, false
}

type RedEyeConfiguration struct {
	API string `json:"API"`
}

func (p *RedEye) Create(configurationRaw json.RawMessage) bool {
	var configuration *RedEyeConfiguration

	configuration = &RedEyeConfiguration{}

	json.Unmarshal(configurationRaw, configuration)
	p.API = &API{Base: configuration.API}

	go p.LoadData()
	return true
}
func (p *RedEye) TypeString() string {
	return "RedEye"
}
func (p *RedEye) IDString() string {
	return p.API.Base
}

type Device struct {
}

func (d *Device) HasChildDevice(find apimodels.Match) bool {
	return false
}
func (d *Device) GetChildDevice(find apimodels.Match) apimodels.Device {
	return &Device{}
}
func (d *Device) IsHidden() bool {
	return false
}
func (d *Device) ListChildren() []apimodels.Device {
	return []apimodels.Device{}
}
func (d *Device) ListActions() []apimodels.Action {
	return []apimodels.Action{}
}
func (d *Device) ListTypes() []apimodels.Type {
	return []apimodels.Type{}
}
func (d *Device) SetValue(float64) bool {
	return false
}
func (d *Device) InvokeAction(string) bool {
	return false
}
func (d *Device) GetName() string {
	return "DeviceName"
}
func (d *Device) GetLocationOne() string {
	return "Floor"
}
func (d *Device) ProviderIDString() string {
	return "RedEyeID"
}
func (d *Device) GetLocationTwo() string {
	return "Room"
}
func (d *Device) IDString() string {
	return "DeviceID"
}
func (d *Device) Matches(apimodels.Match) bool {
	return false
}

type DeviceList struct {
	Rooms []*RedEyeRoom
}

func (l *DeviceList) SetValue(float64) bool {
	return false
}
func (l *DeviceList) List() []apimodels.Device {
	r := make([]apimodels.Device, 0)
	for _, c := range l.Rooms {
		r = append(r, c)
	}
	return r
}
func (l *DeviceList) Add(room *RedEyeRoom) {
	l.Rooms = append(l.Rooms, room)
}
func (l *DeviceList) InvokeAction(string) bool {
	return false
}
