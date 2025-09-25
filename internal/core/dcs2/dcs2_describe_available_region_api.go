package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2DescribeAvailableRegionApi
/* 查询分布式缓存Redis实例支持的所有资源池。
 */type Dcs2DescribeAvailableRegionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeAvailableRegionApi(client *core.CtyunClient) *Dcs2DescribeAvailableRegionApi {
	return &Dcs2DescribeAvailableRegionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/region/getAvailableRegion",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeAvailableRegionApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeAvailableRegionRequest) (*Dcs2DescribeAvailableRegionResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.ResPoolCode != "" {
		ctReq.AddParam("resPoolCode", req.ResPoolCode)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeAvailableRegionResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeAvailableRegionRequest struct {
	RegionId    string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ResPoolCode string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
}

type Dcs2DescribeAvailableRegionResponse struct {
	StatusCode int32                                           `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                          `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  []*Dcs2DescribeAvailableRegionReturnObjResponse `json:"returnObj"`            /*  响应对象  */
	RequestId  string                                          `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                          `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                          `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DescribeAvailableRegionReturnObjResponse struct {
	ResPoolCode string                                                  `json:"resPoolCode,omitempty"` /*  资源池ID  */
	ResPoolName string                                                  `json:"resPoolName,omitempty"` /*  资源池名称  */
	Products    []*Dcs2DescribeAvailableRegionReturnObjProductsResponse `json:"products"`              /*  产品实体  */
}

type Dcs2DescribeAvailableRegionReturnObjProductsResponse struct {
	ProdName      string `json:"prodName,omitempty"`      /*  产品名称  */
	ProdCode      string `json:"prodCode,omitempty"`      /*  产品编码  */
	OuterProdCode string `json:"outerProdCode,omitempty"` /*  产品外部编码  */
	ProdStatus    string `json:"prodStatus,omitempty"`    /*  产品状态<li>2：正常<li>其他表示异常  */
}
