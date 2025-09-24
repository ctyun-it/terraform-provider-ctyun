package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2DescribeHistoryMonitorItemsApi
/* 查询分布式缓存Redis实例或节点支持的性能监控监控项列表。
 */type Dcs2DescribeHistoryMonitorItemsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeHistoryMonitorItemsApi(client *core.CtyunClient) *Dcs2DescribeHistoryMonitorItemsApi {
	return &Dcs2DescribeHistoryMonitorItemsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/resourceMonitor/describeHistoryMonitorItems",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeHistoryMonitorItemsApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeHistoryMonitorItemsRequest) (*Dcs2DescribeHistoryMonitorItemsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeHistoryMonitorItemsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeHistoryMonitorItemsRequest struct {
	RegionId string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
}

type Dcs2DescribeHistoryMonitorItemsResponse struct {
	StatusCode int32                                             `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                            `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2DescribeHistoryMonitorItemsReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                            `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                            `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                            `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DescribeHistoryMonitorItemsReturnObjResponse struct{}
