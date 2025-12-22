package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseListNamespaceV2P2Api
/* 查询指定集群下的namespace列表
 */type CcseListNamespaceV2P2Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseListNamespaceV2P2Api(client *core.CtyunClient) *CcseListNamespaceV2P2Api {
	return &CcseListNamespaceV2P2Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/clusters/{clusterName}/api/v1/namespaces",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseListNamespaceV2P2Api) Do(ctx context.Context, credential core.Credential, req *CcseListNamespaceV2P2Request) (*CcseListNamespaceV2P2Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterName", req.ClusterName)
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
	var resp CcseListNamespaceV2P2Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseListNamespaceV2P2Request struct {
	ClusterName   string `json:"clusterName,omitempty"`   /*  集群id，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
	RegionId      string `json:"regionId,omitempty"`      /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a><br>另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池<br>获取：<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
	LabelSelector string `json:"labelSelector,omitempty"` /*  Kubernetes labelSelector，可通过label过滤资源；label之间通过“,”分隔，特殊符号要转义为url编码，如“=”写为“%3D”
	如：
	{"name":"test"}  */
	FieldSelector string `json:"fieldSelector,omitempty"` /*  Kubernetes fieldSelector，可通过field过滤资源；field之间通过“,”分隔，特殊符号要转义为url编码，如“=”写为“%3D”
	如：
	{"metadata.namespace":"default"}  */
}

type CcseListNamespaceV2P2Response struct {
	StatusCode int32  `json:"statusCode"` /*  响应状态码  */
	Message    string `json:"message"`    /*  响应信息  */
	ReturnObj  string `json:"returnObj"`  /*  返回结果  */
	Error      string `json:"error"`      /*  错误码，参见错误码说明  */
}
