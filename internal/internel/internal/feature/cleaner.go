package feature

import (
	"context"

	"github.com/dronestock/ssl/internal/internel/internal/core"
)

type Cleaner interface {
	Invalidates(ctx *context.Context, certificate *core.Certificate) (certificates []*core.ServerCertificate, err error)

	Delete(ctx *context.Context, certificate *core.ServerCertificate) (err error)
}
