package main

import (
	"context"
)

type stepCleanup struct {
	*plugin
}

func newStepCleanup(plugin *plugin) *stepCleanup {
	return &stepCleanup{
		plugin: plugin,
	}
}

func (r *stepCleanup) Runnable() (runnable bool) {
	if nil != r.Manufacturer.Chuangcache {
		runnable = true
	} else {
		runnable = r.runnable()
	}

	return
}

func (r *stepCleanup) Run(ctx context.Context) (err error) {
	for _, certificate := range r.Certificates {
		err = certificate.Clean(ctx, r.Base)
	}

	return
}

func (r *stepCleanup) runnable() bool {
	return 0 != len(r.Certificates)
}
