package mysql

import ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"

const (
	EndpointNameCtdas = "mysql"
	UrlPordCtdas      = "rds2-global.ctapi.ctyun.cn"
	UrlTestCtdas      = ""
)

var EndpointCtdasTest = ctyunsdk.Endpoint{
	EndpointName: EndpointNameCtdas,
	Url:          UrlTestCtdas,
}

var EndPointCtdasProd = ctyunsdk.Endpoint{
	EndpointName: EndpointNameCtdas,
	Url:          UrlPordCtdas,
}
