package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2GetZonesApi
/* 查询指定资源池下，分布式缓存Redis支持的可用区。
 */type Dcs2GetZonesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2GetZonesApi(client *core.CtyunClient) *Dcs2GetZonesApi {
	return &Dcs2GetZonesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/lifeCycleServant/getZones",
			ContentType:  "",
		},
	}
}

func (a *Dcs2GetZonesApi) Do(ctx context.Context, credential core.Credential, req *Dcs2GetZonesRequest) (*Dcs2GetZonesResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionIdHeader)
	ctReq.AddParam("regionId", req.RegionIdParam)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2GetZonesResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2GetZonesRequest struct {
	RegionIdHeader string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	RegionIdParam  string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
}

type Dcs2GetZonesResponse struct {
	StatusCode int32                          `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                         `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2GetZonesReturnObjResponse `json:"returnObj"`            /*  响应对象  */
	RequestId  string                         `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                         `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                         `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2GetZonesReturnObjResponse struct {
	ZoneList []*Dcs2GetZonesReturnObjZoneListResponse `json:"zoneList"` /*  可用区列表  */
}

type Dcs2GetZonesReturnObjZoneListResponse struct {
	Name          string `json:"name,omitempty"`          /*  可用区名  */
	AzDisplayName string `json:"azDisplayName,omitempty"` /*  可用区展示名称  */
}
