package apimodels

type Match map[string]interface{}

type ControlRequest struct {
	Type   string  `json:"type"`
	Match  Match   `json:"match"`
	Value  float64 `json:"value,omitempty"`
	Action string  `json:"action,omitempty"`
}

type Device interface {
	SetValue(float64) bool
	InvokeAction(string) bool
	ProviderIDString() string
	IDString() string
	GetName() string
	Matches(Match) bool
	ListChildren() []Device
	ListActions() []Action
	ListTypes() []Type
	GetLocationOne() string
	GetLocationTwo() string
	GetChildDevice(Match) Device
	HasChildDevice(Match) bool
	IsHidden() bool
}

type Action interface {
	GetName() string
}

type Type interface {
	GetName() string
}

type Devices interface {
	SetValue(float64) bool
	InvokeAction(string) bool
	List() []Device
}
