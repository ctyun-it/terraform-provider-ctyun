package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseListServicesApi
/* 获取指定集群下的Service列表
 */type CcseListServicesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseListServicesApi(client *core.CtyunClient) *CcseListServicesApi {
	return &CcseListServicesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/clusters/{clusterId}/api/v1/namespaces/{namespaceName}/services",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseListServicesApi) Do(ctx context.Context, credential core.Credential, req *CcseListServicesRequest) (*CcseListServicesResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder = builder.ReplaceUrl("namespaceName", req.NamespaceName)
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
	var resp CcseListServicesResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseListServicesRequest struct {
	ClusterId     string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105" target="_blank">如何获取接口URI中参数</a>  */
	NamespaceName string /*  命名空间名称  */
	RegionId      string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
	LabelSelector string /*  Kubernetes labelSelector，可通过label过滤资源；label之间通过“,”分隔，特殊符号要转义为url编码，如“=”写为“%3D”    */
	FieldSelector string /*  Kubernetes fieldSelector，可通过field过滤资源；label之间通过“,”分隔，特殊符号要转义为url编码，如“=”写为“%3D”  */
}

type CcseListServicesResponse struct {
	StatusCode int32  `json:"statusCode,omitempty"` /*  响应状态码  */
	Message    string `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  string `json:"returnObj,omitempty"`  /*  返回结果  */
	Error      string `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}
