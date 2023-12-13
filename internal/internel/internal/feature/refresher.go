package feature

import (
	"context"

	"github.com/dronestock/ssl/internal/internel/internal/core"
)

type Refresher interface {
	Upload(ctx context.Context, local *core.Certificate) (certificate *core.ServerCertificate, err error)

	Domains(ctx context.Context) (domains []*core.Domain, err error)

	Bind(ctx context.Context, certificate *core.ServerCertificate, domain *core.Domain) (record *core.Record, err error)

	Check(ctx context.Context, record *core.Record) (checked bool, err error)
}
