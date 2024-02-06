package step

import (
	"context"

	"github.com/dronestock/drone"
	"github.com/dronestock/ssl/internal/internel/config"
)

type Cleanup struct {
	base         *drone.Base
	manufacturer *config.Manufacturer
	certificates []*config.Certificate
}

func NewCleanup(base *drone.Base, manufacturer *config.Manufacturer, certificates []*config.Certificate) *Cleanup {
	return &Cleanup{
		base:         base,
		manufacturer: manufacturer,
		certificates: certificates,
	}
}

func (c *Cleanup) Runnable() (runnable bool) {
	if nil != c.manufacturer.Chuangcache {
		runnable = true
	} else {
		runnable = c.runnable()
	}

	return
}

func (c *Cleanup) Run(ctx *context.Context) (err error) {
	for _, certificate := range c.certificates {
		err = certificate.Clean(ctx, c.base, &certificate.Certificate)
	}

	return
}

func (c *Cleanup) runnable() bool {
	return 0 != len(c.certificates)
}
