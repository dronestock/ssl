package chuangcache

type Certificate struct {
	Title  string            `json:"Title,omitempty"`
	Key    string            `json:"CertKey,omitempty"`
	Status CertificateStatus `json:"Status,omitempty"`
}
