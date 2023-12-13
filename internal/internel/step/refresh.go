package step

import (
	"context"

	"github.com/dronestock/drone"
	"github.com/dronestock/ssl/internal/internel/config"
)

type Refresh struct {
	base         *drone.Base
	manufacturer *config.Manufacturer
	certificates []*config.Certificate
}

func NewRefresh(base *drone.Base, manufacturer *config.Manufacturer, certificates []*config.Certificate) *Refresh {
	return &Refresh{
		base:         base,
		manufacturer: manufacturer,
		certificates: certificates,
	}
}

func (r *Refresh) Runnable() (runnable bool) {
	if nil != r.manufacturer.Chuangcache {
		runnable = true
	} else {
		runnable = r.runnable()
	}

	return
}

func (r *Refresh) Run(ctx context.Context) (err error) {
	for _, certificate := range r.certificates {
		err = certificate.Refresh(ctx, r.base, certificate)
	}

	return
}

func (r *Refresh) runnable() bool {
	return 0 != len(r.certificates)
}
