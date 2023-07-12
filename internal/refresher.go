package internal

import (
	"context"
)

type Refresher interface {
	Init(ctx context.Context) (err error)

	Upload(ctx context.Context, certificate *Certificate) (id string, err error)

	Domains(ctx context.Context) (domains []*Domain, err error)

	Bind(ctx context.Context, id string, domain *Domain) (err error)
}
