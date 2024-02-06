package apisix

type Value struct {
	Id   string   `json:"id,omitempty"`
	SNIs []string `json:"snis,omitempty"`
	Key  string   `json:"key,omitempty"`
	Cert string   `json:"cert,omitempty"`
}
