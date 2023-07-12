package main

import (
	"context"
)

type stepTencent struct {
	*plugin
}

func newStepTencent(plugin *plugin) *stepTencent {
	return &stepTencent{
		plugin: plugin,
	}
}

func (st *stepTencent) Runnable() (runnable bool) {
	if nil != st.Manufacturer.Tencent {
		runnable = true
	} else {
		runnable = st.runnable()
	}

	return
}

func (st *stepTencent) Run(ctx context.Context) (err error) {
	for _, _certificate := range st.Certificates {
		if nil != _certificate.Tencent {
			err = _certificate.Tencent.refresh(ctx, st.Base, _certificate)
		} else if nil != st.Tencent {
			err = st.Tencent.refresh(ctx, st.Base, _certificate)
		}
	}

	return
}

func (st *stepTencent) runnable() (runnable bool) {
	for _, _certificate := range st.Certificates {
		if nil != _certificate.Tencent {
			runnable = true
		}
		if runnable {
			break
		}
	}

	return
}
