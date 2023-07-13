package config

import (
	"github.com/dronestock/ssl/internal/core"
)

type Certificate struct {
	Manufacturer
	core.Certificate

	// 环境变量
	Environments map[string]string `json:"environments,omitempty"`
}
