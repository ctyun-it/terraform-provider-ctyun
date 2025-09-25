package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2QueryAnalysisInstanceTasksInfoApi
/* 查询分布式缓存Redis诊断分析报告详情，内容包含实例可用性、数据同步、负载状态、存储、网络、慢请求、命令耗时统计等信息。
 */type Dcs2QueryAnalysisInstanceTasksInfoApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2QueryAnalysisInstanceTasksInfoApi(client *core.CtyunClient) *Dcs2QueryAnalysisInstanceTasksInfoApi {
	return &Dcs2QueryAnalysisInstanceTasksInfoApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/keyAnalysisMgrServant/queryAnalysisInstanceTasksInfo",
			ContentType:  "",
		},
	}
}

func (a *Dcs2QueryAnalysisInstanceTasksInfoApi) Do(ctx context.Context, credential core.Credential, req *Dcs2QueryAnalysisInstanceTasksInfoRequest) (*Dcs2QueryAnalysisInstanceTasksInfoResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("taskId", req.TaskId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2QueryAnalysisInstanceTasksInfoResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2QueryAnalysisInstanceTasksInfoRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例ID  */
	TaskId     string /*  任务ID  */
}

type Dcs2QueryAnalysisInstanceTasksInfoResponse struct {
	StatusCode int32  `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string `json:"message,omitempty"`    /*  响应信息  */
	RequestId  string `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string `json:"code,omitempty"`       /*  响应码描述  */
	Error      string `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}
