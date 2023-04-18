package main

import (
	"context"
)

type stepChuangcache struct {
	*plugin
}

func newStepChuangcache(plugin *plugin) *stepChuangcache {
	return &stepChuangcache{
		plugin: plugin,
	}
}

func (sc *stepChuangcache) Runnable() (runnable bool) {
	if nil != sc.Manufacturer.Chuangcache {
		runnable = true
	} else {
		runnable = sc.runnable()
	}

	return
}

func (sc *stepChuangcache) Run(ctx context.Context) (err error) {
	for _, _certificate := range sc.Certificates {
		if nil != _certificate.Chuangcache {
			err = _certificate.Chuangcache.refresh(ctx, sc.Base, _certificate)
		} else if nil != sc.Chuangcache {
			err = sc.Chuangcache.refresh(ctx, sc.Base, _certificate)
		}
	}

	return
}

func (sc *stepChuangcache) runnable() (runnable bool) {
	for _, _certificate := range sc.Certificates {
		if nil != _certificate.Chuangcache {
			runnable = true
		}
		if runnable {
			break
		}
	}

	return
}
