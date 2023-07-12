package refresher

import (
	"context"
	"fmt"

	"github.com/dronestock/ssl/internal"
	"github.com/dronestock/ssl/internal/chuangcache"
	"github.com/dronestock/ssl/internal/config"
	"github.com/dronestock/ssl/internal/tencent"
	"github.com/goexl/gox/field"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

var _ internal.Refresher = (*Tencent)(nil)

type Tencent struct {
	id  string
	key string

	ssl    *ssl.Client
	cdn    *cdn.Client
	titles map[string]bool
}

func NewTencent(config *config.Tencent) *Tencent {
	return &Tencent{
		id:  config.Id,
		key: config.Key,
	}
}

func (t *Tencent) Init(_ context.Context) (err error) {
	credential := common.NewCredential(t.id, t.key)
	if sc, se := ssl.NewClient(credential, regions.Chengdu, profile.NewClientProfile()); nil != se {
		err = se
	} else {
		t.ssl = sc
	}

	return
}

func (t *Tencent) Upload(_ context.Context, certificate *internal.Certificate) (id string, err error) {
	req := new(tencent.UploadReq)
	req.Alias = common.StringPtr(certificate.Title)
	if le := certificate.Load(req); nil != le {
		err = le
	} else if rsp, uce := t.ssl.UploadCertificate(req.Request()); nil != uce {
		err = uce
	} else {
		id = *rsp.Response.CertificateId
		// t.titles[req.Title] = true
	}

	return
}

func (t *Tencent) Bind(ctx context.Context, id string, domain *internal.Domain) (err error) {
	req := new(ssl.DeployCertificateInstanceRequest)
	req.CertificateId = common.StringPtr(id)
	req.InstanceIdList = []*string{common.StringPtr(domain.Id)}
	req.ResourceType = domain.Type
	rsp := new(chuangcache.Response[bool])
	url := fmt.Sprintf("%s/%s", chuangcache.ApiEndpoint, "config/bindDomainCertificate")
	if ce := t.call(ctx, base, url, req, rsp); nil != ce {
		err = ce
	} else if !rsp.Data {
		err = ext.NewFields("绑定证书失败", field.New("req", req), field.New("rsp", rsp))
	}

	return
}

func (t *Tencent) Domains(ctx context.Context) (domains []*internal.Domain, err error) {
	domains = make([]*internal.Domain, 0, 1)
	if ce := t.cdnDomains(ctx, &domains); nil != ce {
		err = ce
	}

	return
}

func (t *Tencent) cdnDomains(_ context.Context, domains *[]*internal.Domain) (err error) {
	req := new(cdn.DescribeDomainsRequest)
	req.Limit = common.Int64Ptr(1000)
	if rsp, dde := t.cdn.DescribeDomains(req); nil != dde {
		err = dde
	} else if 0 != len(rsp.Response.Domains) {
		for _, brief := range rsp.Response.Domains {
			domain := new(internal.Domain)
			domain.Id = *brief.ResourceId
			domain.Name = *brief.Domain
			domain.Type = internal.DomainTypeCdn

			*domains = append(*domains, domain)
		}
	}

	return
}
