package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseHasTemplateInstanceExistedApi
/* 调用该接口可检查模板实例是否存在。
 */type CcseHasTemplateInstanceExistedApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseHasTemplateInstanceExistedApi(client *core.CtyunClient) *CcseHasTemplateInstanceExistedApi {
	return &CcseHasTemplateInstanceExistedApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/clusters/{clusterId}/namespaces/{namespaceName}/templateinstance/{templateInstanceName}/exists",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseHasTemplateInstanceExistedApi) Do(ctx context.Context, credential core.Credential, req *CcseHasTemplateInstanceExistedRequest) (*CcseHasTemplateInstanceExistedResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder = builder.ReplaceUrl("namespaceName", req.NamespaceName)
	builder = builder.ReplaceUrl("templateInstanceName", req.TemplateInstanceName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseHasTemplateInstanceExistedResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseHasTemplateInstanceExistedRequest struct {
	ClusterId            string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
	NamespaceName        string /*  命名空间名称  */
	TemplateInstanceName string /*  模板实例名称  */
	RegionId             string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
}

type CcseHasTemplateInstanceExistedResponse struct {
	StatusCode int32  `json:"statusCode,omitempty"` /*  状态码  */
	RequestId  string `json:"requestId,omitempty"`  /*  请求id  */
	Message    string `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  *bool  `json:"returnObj"`            /*  返回对象  */
	Error      string `json:"error,omitempty"`      /*  错误码  */
}
