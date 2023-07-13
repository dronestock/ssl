package feature

import (
	"context"

	"github.com/dronestock/ssl/internal/core"
)

type Cleaner interface {
	Invalidates(ctx context.Context) (certificates []*core.ServerCertificate, err error)

	Delete(ctx context.Context, certificate *core.ServerCertificate) (deleted bool, err error)
}
