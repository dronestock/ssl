package main

import (
	"context"
)

type stepRefresh struct {
	*plugin
}

func newStepRefresh(plugin *plugin) *stepRefresh {
	return &stepRefresh{
		plugin: plugin,
	}
}

func (r *stepRefresh) Runnable() (runnable bool) {
	if nil != r.Manufacturer.Chuangcache {
		runnable = true
	} else {
		runnable = r.runnable()
	}

	return
}

func (r *stepRefresh) Run(ctx context.Context) (err error) {
	for _, certificate := range r.Certificates {
		err = certificate.Refresh(ctx, r.Base, certificate)
	}

	return
}

func (r *stepRefresh) runnable() bool {
	return 0 != len(r.Certificates)
}
