package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseListNodesApi
/* 获取指定集群下的node列表
 */type CcseListNodesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseListNodesApi(client *core.CtyunClient) *CcseListNodesApi {
	return &CcseListNodesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/clusters/{clusterId}/api/v1/nodes",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseListNodesApi) Do(ctx context.Context, credential core.Credential, req *CcseListNodesRequest) (*CcseListNodesResponse, error) {
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
	var resp CcseListNodesResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseListNodesRequest struct {
	ClusterId string `json:"clusterId,omitempty"` /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
	RegionId  string `json:"regionId,omitempty"`  /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
	LabelSelector string `json:"labelSelector,omitempty"` /*  Kubernetes labelSelector，可通过label过滤资源；label之间通过“,”分隔，特殊符号要转义为url编码，如“=”写为“%3D”   */
	FieldSelector string `json:"fieldSelector,omitempty"` /*  Kubernetes fieldSelector，可通过field过滤资源；label之间通过“,”分隔，特殊符号要转义为url编码，如“=”写为“%3D”  */
}

type CcseListNodesResponse struct {
	StatusCode int32  `json:"statusCode"` /*  响应状态码  */
	Message    string `json:"message"`    /*  响应信息  */
	ReturnObj  string `json:"returnObj"`  /*  返回结果  */
	Error      string `json:"error"`      /*  错误码，参见错误码说明  */
}
