package manufacturer

import (
	"context"

	"github.com/dronestock/ssl/internal/core"
	"github.com/dronestock/ssl/internal/feature"
	"github.com/dronestock/ssl/internal/tencent"
	"github.com/go-resty/resty/v2"
	"github.com/goexl/exc"
	"github.com/goexl/gox/field"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

var (
	_ feature.Refresher = (*Tencent)(nil)
	_ feature.Cleaner   = (*Tencent)(nil)
)

type Tencent struct {
	config *core.Tencent

	http *resty.Client
	ssl  *ssl.Client
	cdn  *cdn.Client
}

func NewTencent(config *core.Tencent) *Tencent {
	return &Tencent{
		config: config,
	}
}

func (t *Tencent) Init(_ context.Context) (err error) {
	credential := common.NewCredential(t.config.Id, t.config.Key)
	cp := profile.NewClientProfile()
	if sc, se := ssl.NewClient(credential, t.config.Region, cp); nil != se {
		err = se
	} else {
		t.ssl = sc
	}

	return
}

func (t *Tencent) Upload(ctx context.Context, certificate *core.Certificate) (id string, err error) {
	req := new(tencent.UploadReq)
	req.Alias = common.StringPtr(certificate.Title)
	if le := certificate.Load(req); nil != le {
		err = le
	} else if rsp, uce := t.ssl.UploadCertificateWithContext(ctx, req.Request()); nil != uce {
		err = uce
	} else {
		id = *rsp.Response.CertificateId
	}

	return
}

func (t *Tencent) Bind(ctx context.Context, id string, domain *core.Domain) (err error) {
	req := new(ssl.DeployCertificateInstanceRequest)
	req.CertificateId = common.StringPtr(id)
	req.InstanceIdList = []*string{common.StringPtr(domain.Id)}
	req.ResourceType = domain.TencentType()
	if rsp, dce := t.ssl.DeployCertificateInstanceWithContext(ctx, req); nil != dce {
		err = dce
	} else if 0 == *rsp.Response.DeployStatus {
		err = exc.NewFields("绑定证书失败", field.New("req", req), field.New("rsp", rsp))
	}

	return
}

func (t *Tencent) Domains(ctx context.Context) (domains []*core.Domain, err error) {
	domains = make([]*core.Domain, 0, 1)
	if ce := t.cdnDomains(ctx, &domains); nil != ce {
		err = ce
	}

	return
}

func (t *Tencent) Certificates(ctx context.Context) (certificates []*core.ServerCertificate, err error) {
	return
}

func (t *Tencent) Delete(ctx context.Context, certificate *core.ServerCertificate) (err error) {
	return
}

func (t *Tencent) cdnDomains(ctx context.Context, domains *[]*core.Domain) (err error) {
	req := new(cdn.DescribeDomainsRequest)
	req.Limit = common.Int64Ptr(1000)
	if rsp, dde := t.cdn.DescribeDomainsWithContext(ctx, req); nil != dde {
		err = dde
	} else if 0 != len(rsp.Response.Domains) {
		for _, brief := range rsp.Response.Domains {
			domain := new(core.Domain)
			domain.Id = *brief.ResourceId
			domain.Name = *brief.Domain
			domain.Type = core.DomainTypeCdn

			*domains = append(*domains, domain)
		}
	}

	return
}
