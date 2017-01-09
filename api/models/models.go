package apimodels

type Match map[string]interface{}

type ControlRequest struct {
	Match Match   `json:"match"`
	Value float64 `json:"value"`
}
