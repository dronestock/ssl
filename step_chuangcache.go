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

func (sc *stepChuangcache) Runnable() bool {
	return true
}

func (sc *stepChuangcache) Run(_ context.Context) (err error) {
	return
}
