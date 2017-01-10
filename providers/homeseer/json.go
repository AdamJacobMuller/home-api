package homeseer

type JD_HSControlDevices struct {
	Version string                `json:"Version"`
	Name    string                `json:"Name"`
	Devices []*JD_HSControlDevice `json:"Devices"`
}
type JD_HSControlDevice struct {
	ID           int                 `json:"ref"`
	ControlPairs []*JD_HSControlPair `json:"ControlPairs"`
}
type JD_HSControlRange struct {
	Start        int    `json:"RangeStart"`
	End          int    `json:"RangeEnd"`
	StatusPrefix string `json:"RangeStatusPrefix"`
	StatusSuffix string `json:"RangeStatusSuffix"`
}
type JD_HSControlPair struct {
	Label    string             `json:"Label"`
	DeviceID int                `json:"ref"`
	Type     int                `json:"ControlType"`
	Value    float64            `json:"ControlValue"`
	Range    *JD_HSControlRange `json:"Range"`
}

func (h *JD_HSControlPair) asHS() *HSAction {
	return &HSAction{
		Label:    h.Label,
		DeviceID: h.DeviceID,
		Value:    h.Value,
	}
}

type JD_HSDevices struct {
	Version string         `json:"Version"`
	Name    string         `json:"Name"`
	Devices []*JD_HSDevice `json:"Devices"`
}
type JD_HSDevice struct {
	ID          int     `json:"ref"`
	Name        string  `json:"name"`
	LocationOne string  `json:"location"`
	LocationTwo string  `json:"location2"`
	Value       float64 `json:"value"`
	Status      string  `json:"status"`
	LastChange  string  `json:"last_change"`
	TypeString  string  `json:"device_type_string"`

	Relationship       int   `json:"relationship"`
	Associated_devices []int `json:"associated_devices"`
}

func (h *JD_HSDevice) asHS() *HSDevice {
	return &HSDevice{
		Name:        h.Name,
		ID:          h.ID,
		LocationOne: h.LocationOne,
		LocationTwo: h.LocationTwo,
		Value:       h.Value,
		Status:      h.Status,
		TypeString:  h.TypeString,
	}
}
