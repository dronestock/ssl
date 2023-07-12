package main

import (
	"github.com/dronestock/ssl/internal/config"
)

type Manufacturer struct {
	Chuangcache *config.Chuangcache `default:"${CHUANGCACHE}" json:"chuangcache,omitempty"`
	Tencent     *config.Tencent     `default:"${TENCENT}" json:"tencent,omitempty"`
}
