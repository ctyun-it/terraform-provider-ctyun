package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2DescribeInstanceHistoryMonitorValuesApi
/* 查询分布式缓存Redis实例历史性能监控数据。
 */type Dcs2DescribeInstanceHistoryMonitorValuesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeInstanceHistoryMonitorValuesApi(client *core.CtyunClient) *Dcs2DescribeInstanceHistoryMonitorValuesApi {
	return &Dcs2DescribeInstanceHistoryMonitorValuesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/resourceMonitor/describeInstanceHistoryMonitorValues",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeInstanceHistoryMonitorValuesApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeInstanceHistoryMonitorValuesRequest) (*Dcs2DescribeInstanceHistoryMonitorValuesResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("startTime", req.StartTime)
	ctReq.AddParam("endTime", req.EndTime)
	ctReq.AddParam("type", req.RawType)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeInstanceHistoryMonitorValuesResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeInstanceHistoryMonitorValuesRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例ID  */
	StartTime  string /*  开始时间  */
	EndTime    string /*  结束时间  */
	RawType    string /*  实例监控类型<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=15090&isNormal=1&vid=270">查询性能监控指标列表</a> instanceMonitorList列表type字段  */
}

type Dcs2DescribeInstanceHistoryMonitorValuesResponse struct {
	StatusCode int32                                                      `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                                     `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2DescribeInstanceHistoryMonitorValuesReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                                     `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                                     `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                                     `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DescribeInstanceHistoryMonitorValuesReturnObjResponse struct{}
