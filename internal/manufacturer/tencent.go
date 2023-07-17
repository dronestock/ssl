package manufacturer

import (
	"context"
	"math"

	"github.com/dronestock/ssl/internal/core"
	"github.com/dronestock/ssl/internal/feature"
	"github.com/dronestock/ssl/internal/tencent"
	"github.com/go-resty/resty/v2"
	"github.com/goexl/exc"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/goexl/simaqian"
	api "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

var (
	_ feature.Refresher = (*Tencent)(nil)
	_ feature.Cleaner   = (*Tencent)(nil)
)

type Tencent struct {
	simaqian.Logger

	http *resty.Client
	ssl  *ssl.Client
	cdn  *cdn.Client
	api  *api.Client
}

func NewTencent(config *core.Tencent, logger simaqian.Logger) (tencent *Tencent, err error) {
	tencent = new(Tencent)
	tencent.Logger = logger

	credential := common.NewCredential(config.Id, config.Key)
	cp := profile.NewClientProfile()
	if sc, se := ssl.NewClient(credential, regions.Chengdu, cp); nil != se {
		err = se
	} else if cc, ce := cdn.NewClient(credential, regions.Chengdu, cp); nil != ce {
		err = ce
	} else if ac, ae := api.NewClient(credential, regions.Chengdu, cp); nil != ae {
		err = ae
	} else {
		tencent.ssl = sc
		tencent.cdn = cc
		tencent.api = ac
	}

	return
}

func (t *Tencent) Upload(ctx context.Context, local *core.Certificate) (cert *core.ServerCertificate, err error) {
	req := tencent.NewUploadReq()
	req.Alias(local.Title)
	if le := local.Load(req); nil != le {
		err = le
	} else if rsp, uce := t.ssl.UploadCertificateWithContext(ctx, req.Request()); nil != uce {
		err = uce
	} else {
		cert = new(core.ServerCertificate)
		cert.Id = *rsp.Response.CertificateId
		cert.Title = local.Title
	}

	return
}

func (t *Tencent) Bind(
	ctx context.Context,
	cert *core.ServerCertificate,
	domain *core.Domain,
) (record *core.Record, err error) {
	req := ssl.NewDeployCertificateInstanceRequest()
	req.CertificateId = common.StringPtr(cert.Id)
	req.InstanceIdList = []*string{common.StringPtr(domain.Name)}
	req.ResourceType = domain.TencentType()
	if rsp, dce := t.ssl.DeployCertificateInstanceWithContext(ctx, req); nil != dce {
		err = dce
	} else if 0 == *rsp.Response.DeployStatus {
		err = exc.NewFields("绑定证书失败", field.New("req", req), field.New("rsp", rsp))
	} else {
		record = new(core.Record)
		record.Id = gox.ToString(*rsp.Response.DeployRecordId)
	}

	return
}

func (t *Tencent) Check(ctx context.Context, record *core.Record) (checked bool, err error) {
	req := ssl.NewDescribeHostDeployRecordDetailRequest()
	req.DeployRecordId = common.StringPtr(record.Id)
	if rsp, dce := t.ssl.DescribeHostDeployRecordDetailWithContext(ctx, req); nil != dce {
		err = dce
	} else {
		checked = t.checkDeploy(rsp.Response.DeployRecordDetailList)
	}

	return
}

func (t *Tencent) Domains(ctx context.Context) (domains []*core.Domain, err error) {
	domains = make([]*core.Domain, 0, 1)
	if ce := t.cdnDomains(ctx, &domains); nil != ce {
		err = ce
	} else if ae := t.apiDomains(ctx, &domains); nil != ae {
		err = ae
	}

	return
}

func (t *Tencent) Invalidates(ctx context.Context) (certificates []*core.ServerCertificate, err error) {
	certificates = make([]*core.ServerCertificate, 0, 1)
	req := ssl.NewDescribeCertificatesRequest()
	req.Deployable = common.Uint64Ptr(1)
	req.Limit = common.Uint64Ptr(1000)
	for page := 0; page < math.MaxInt; page++ {
		req.Offset = common.Uint64Ptr(uint64(page) * (*req.Limit))
		if rsp, dce := t.ssl.DescribeCertificatesWithContext(ctx, req); nil != dce {
			err = dce
		} else if 0 == len(rsp.Response.Certificates) {
			break
		} else {
			for _, cert := range rsp.Response.Certificates {
				certificate := new(core.ServerCertificate)
				certificate.Id = *cert.CertificateId
				certificate.Title = *cert.Alias

				certificates = gox.Ift(t.invalidate(ctx, cert), append(certificates, certificate), certificates)
			}
		}
	}

	return
}

func (t *Tencent) Delete(ctx context.Context, cert *core.ServerCertificate) (deleted bool, err error) {
	req := ssl.NewDeleteCertificateRequest()
	req.CertificateId = common.StringPtr(cert.Id)
	if rsp, dce := t.ssl.DeleteCertificateWithContext(ctx, req); nil != dce {
		err = dce
	} else {
		deleted = *rsp.Response.DeleteResult
	}

	return
}

func (t *Tencent) cdnDomains(ctx context.Context, domains *[]*core.Domain) (err error) {
	req := cdn.NewDescribeDomainsRequest()
	req.Limit = common.Int64Ptr(1000)
	for page := 0; page < math.MaxInt; page++ {
		req.Offset = common.Int64Ptr(int64(page) * (*req.Limit))
		if rsp, dde := t.cdn.DescribeDomainsWithContext(ctx, req); nil != dde {
			err = dde
		} else if 0 == len(rsp.Response.Domains) {
			break
		} else {
			for _, brief := range rsp.Response.Domains {
				domain := new(core.Domain)
				domain.Id = *brief.ResourceId
				domain.Name = *brief.Domain
				domain.Type = core.DomainTypeCdn

				*domains = append(*domains, domain)
			}
		}
	}

	return
}

func (t *Tencent) apiDomains(ctx context.Context, domains *[]*core.Domain) (err error) {
	req := api.NewDescribeServiceSubDomainsRequest()
	req.Limit = common.Int64Ptr(1000)
	for page := 0; page < math.MaxInt; page++ {
		req.Offset = common.Int64Ptr(int64(page) * (*req.Limit))
		if rsp, dce := t.api.DescribeServiceSubDomainsWithContext(ctx, req); nil != dce {
			err = dce
		} else if 0 == len(rsp.Response.Result.DomainSet) {
			break
		} else {
			for _, _domain := range rsp.Response.Result.DomainSet {
				domain := new(core.Domain)
				domain.Id = *_domain.DomainName
				domain.Name = *_domain.DomainName
				domain.Type = core.DomainTypeGateway

				*domains = append(*domains, domain)
			}
		}
	}

	return
}

func (t *Tencent) invalidate(ctx context.Context, certificate *ssl.Certificates) (invalidate bool) {
	total := 0
	req := ssl.NewDescribeDeployedResourcesRequest()
	req.CertificateIds = []*string{certificate.CertificateId}
	fields := gox.Fields[any]{
		field.New("id", certificate.CertificateId),
		field.New("title", certificate.Alias),
	}

	// 检查内容分发网络
	t.Info("检查内容分发网络是否有关联证书", fields...)
	req.ResourceType = common.StringPtr("cdn")
	if rsp, dre := t.ssl.DescribeDeployedResourcesWithContext(ctx, req); nil != dre {
		total += 1
		t.Warn("检查失效证书出错", field.Error(dre))
	} else {
		total += t.total(rsp.Response.DeployedResources)
	}

	// 检查网关
	t.Info("检查负载均衡是否有关联证书", fields...)
	req.ResourceType = common.StringPtr("clb")
	if rsp, dre := t.ssl.DescribeDeployedResourcesWithContext(ctx, req); nil != dre {
		total += 1
		t.Warn("检查失效证书出错", field.Error(dre))
	} else {
		total += t.total(rsp.Response.DeployedResources)
	}

	// 检查内容分发网络
	t.Info("检查云直播是否有关联证书", fields...)
	req.ResourceType = common.StringPtr("live")
	if rsp, dre := t.ssl.DescribeDeployedResourcesWithContext(ctx, req); nil != dre {
		total += 1
		t.Warn("检查失效证书出错", field.Error(dre))
	} else {
		total += t.total(rsp.Response.DeployedResources)
	}

	// 检查内容分发网络
	t.Info("检查网络防火墙是否有关联证书", fields...)
	req.ResourceType = common.StringPtr("waf")
	if rsp, dre := t.ssl.DescribeDeployedResourcesWithContext(ctx, req); nil != dre {
		total += 1
		t.Warn("检查失效证书出错", field.Error(dre))
	} else {
		total += t.total(rsp.Response.DeployedResources)
	}

	// 检查内容分发网络
	t.Info("检查网络攻击是否有关联证书", fields...)
	req.ResourceType = common.StringPtr("antiddos")
	if rsp, dre := t.ssl.DescribeDeployedResourcesWithContext(ctx, req); nil != dre {
		total += 1
		t.Warn("检查失效证书出错", field.Error(dre))
	} else {
		total += t.total(rsp.Response.DeployedResources)
	}

	// 实例数为0表示已失效
	invalidate = 0 == total

	return
}

func (t *Tencent) total(resources []*ssl.DeployedResources) (total int) {
	for _, resource := range resources {
		total += len(resource.Resources)
	}

	return
}

func (t *Tencent) checkDeploy(details []*ssl.DeployRecordDetail) (checked bool) {
	for _, detail := range details {
		if 1 == *detail.Status {
			checked = true
		}

		if checked {
			break
		}
	}

	return
}
