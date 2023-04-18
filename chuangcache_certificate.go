package main

type chuangcacheCertificate struct {
	Title  string                       `json:"Title,omitempty"`
	Key    string                       `json:"CertKey,omitempty"`
	Status chuangcacheCertificateStatus `json:"Status,omitempty"`
}
