package apimodels

type Match map[string]interface{}

type ControlRequest struct {
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
	GetLocationOne() string
	GetLocationTwo() string
	GetTypes() []string
}

type Action interface {
	GetName() string
}

type Devices interface {
	SetValue(float64) bool
	InvokeAction(string) bool
	List() []Device
}
