package main

type chuangcacheCertificate struct {
	Key    string                       `json:"CertKey,omitempty"`
	Status chuangcacheCertificateStatus `json:"Status,omitempty"`
}
