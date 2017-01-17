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
	IDString() string
	Matches(Match) bool
	ListActions() []Action
}

type Action interface {
	GetName() string
}

type Devices interface {
	SetValue(float64) bool
	InvokeAction(string) bool
}
