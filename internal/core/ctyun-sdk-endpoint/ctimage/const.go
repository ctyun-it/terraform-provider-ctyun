package ctimage

import "terraform-provider-ctyun/internal/core/ctyun-sdk-core"

const (
	EndpointNameCtimage = "ctimage"
	UrlProdCtiamge      = "ctimage-global.ctapi.ctyun.cn"
	UrlTestCtiamge      = "ctimage-global.ctapi-test.ctyun.cn:21443"
)

var EndpointCtimageProd = ctyunsdk.Endpoint{
	EndpointName: EndpointNameCtimage,
	Url:          UrlProdCtiamge,
}

var EndpointCtimageTest = ctyunsdk.Endpoint{
	EndpointName: EndpointNameCtimage,
	Url:          UrlTestCtiamge,
}
