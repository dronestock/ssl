package feature

import (
	"context"

	"github.com/dronestock/ssl/internal/core"
)

type Refresher interface {
	Init(ctx context.Context) (err error)

	Upload(ctx context.Context, certificate *core.Certificate) (id string, err error)

	Domains(ctx context.Context) (domains []*core.Domain, err error)

	Bind(ctx context.Context, id string, domain *core.Domain) (err error)
}
