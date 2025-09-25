package ctecs

import "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"

const (
	EndpointNameCtecs = "ctecs"
	UrlProdCtecs      = "ctecs-global.ctapi.ctyun.cn"
	UrlTestCtecs      = "ctecs-global.ctapi-test.ctyun.cn:21443"
)

var EndpointCtecsProd = ctyunsdk.Endpoint{
	EndpointName: EndpointNameCtecs,
	Url:          UrlProdCtecs,
}

var EndpointCtecsTest = ctyunsdk.Endpoint{
	EndpointName: EndpointNameCtecs,
	Url:          UrlTestCtecs,
}
