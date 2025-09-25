package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseListEventsApi
/* 调用该接口查询集群一小时内的事件列表。
 */type CcseListEventsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseListEventsApi(client *core.CtyunClient) *CcseListEventsApi {
	return &CcseListEventsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/clusters/{clusterId}/apis/events.k8s.io/v1/events",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseListEventsApi) Do(ctx context.Context, credential core.Credential, req *CcseListEventsRequest) (*CcseListEventsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.LabelSelector != "" {
		ctReq.AddParam("labelSelector", req.LabelSelector)
	}
	if req.FieldSelector != "" {
		ctReq.AddParam("fieldSelector", req.FieldSelector)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseListEventsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseListEventsRequest struct {
	ClusterId string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
	RegionId  string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&amp;api=5851&amp;data=87">资源池列表查询</a>  */
	LabelSelector string /*  Kubernetes labelSelector，可通过label过滤资源；label之间通过“,”分隔，特殊符号要转义为url编码，如“=”写为“%3D”  */
	FieldSelector string /*  Kubernetes fieldSelector，可通过field过滤资源；label之间通过“,”分隔，特殊符号要转义为url编码，如“=”写为“%3D”，以下为过滤Event常用的key及其含义：
	type：事件等级，有Normal和Warning两种等级
	regarding.kind：产生事件资源的类型，如Pod、Node等
	regarding.name：产生事件资源的名称
	regarding.namesapce：产生事件资源所在的命名空间   */
}

type CcseListEventsResponse struct {
	StatusCode int32  `json:"statusCode,omitempty"` /*  响应状态码  */
	Message    string `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  string `json:"returnObj,omitempty"`  /*  返回结果，yaml格式事件列表  */
	Error      string `json:"error,omitempty"`      /*  错误码  */
}
