package chuangcache

import (
	"github.com/dronestock/ssl/internal/core"
)

type Certificate struct {
	Title  string            `json:"Title,omitempty"`
	Key    string            `json:"CertKey,omitempty"`
	Status CertificateStatus `json:"Status,omitempty"`
}

func (c *Certificate) InternalStatus() (status core.CertificateStatus) {
	switch c.Status {
	case CertificateStatusInuse:
		status = core.CertificateStatusInuse
	}

	return
}
