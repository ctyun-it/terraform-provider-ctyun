package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseCreateClusterDelegateApi
/* 接口将创建CCE服务委托，该委托授权CCE代管云资源。需在创建集群之前先调用该接口完成服务委托创建，否则CCE不被授权，将无法正常管理集群
 */type CcseCreateClusterDelegateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseCreateClusterDelegateApi(client *core.CtyunClient) *CcseCreateClusterDelegateApi {
	return &CcseCreateClusterDelegateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/cce/clusters/delegate/createdelegate",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseCreateClusterDelegateApi) Do(ctx context.Context, credential core.Credential, req *CcseCreateClusterDelegateRequest) (*CcseCreateClusterDelegateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseCreateClusterDelegateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseCreateClusterDelegateRequest struct {
	RegionId string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>
	*/
}

type CcseCreateClusterDelegateResponse struct {
	StatusCode int32  `json:"statusCode,omitempty"` /*  状态码  */
	Message    string `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  *bool  `json:"returnObj"`            /*  返回对象  */
	Error      string `json:"error,omitempty"`      /*  返回错误码  */
}
