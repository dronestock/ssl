package apisix

type Certificate struct {
	Index int   `json:"createdIndex,omitempty"`
	Value Value `json:"value,omitempty"`
}
