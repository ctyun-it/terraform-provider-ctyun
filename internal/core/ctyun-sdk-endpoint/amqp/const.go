package amqp

import ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"

const (
	EndpointName = "amqp"
	UrlProd      = "amqp-global.ctapi.ctyun.cn"
	UrlTest      = ""
)

var EndpointTest = ctyunsdk.Endpoint{
	EndpointName: EndpointName,
	Url:          UrlTest,
}

var EndPointProd = ctyunsdk.Endpoint{
	EndpointName: EndpointName,
	Url:          UrlProd,
}
