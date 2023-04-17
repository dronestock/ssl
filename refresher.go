package main

import (
	"context"

	"github.com/dronestock/drone"
)

type refresher interface {
	refresh(ctx context.Context, base *drone.Base, certificate *certificate) (err error)
}
