package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseGetClusterAutoscalerPolicyApi
/* 查询节点弹性伸缩策略详情
 */type CcseGetClusterAutoscalerPolicyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseGetClusterAutoscalerPolicyApi(client *core.CtyunClient) *CcseGetClusterAutoscalerPolicyApi {
	return &CcseGetClusterAutoscalerPolicyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/clusters/{clusterId}/apis/autoscaler.ccse.ctyun.cn/v1/horizontalnodeautoscalers/{name}",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseGetClusterAutoscalerPolicyApi) Do(ctx context.Context, credential core.Credential, req *CcseGetClusterAutoscalerPolicyRequest) (*CcseGetClusterAutoscalerPolicyResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder = builder.ReplaceUrl("name", req.Name)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseGetClusterAutoscalerPolicyResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseGetClusterAutoscalerPolicyRequest struct {
	ClusterId string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。

	如：
	/v1.2/cce/clusters/{clusterId}/apis/autoscaler.ccse.ctyun.cn/v1/horizontalnodeautoscaler/{name}  */
	Name string /*  节点弹性伸缩策略名称

	如：
	/v1.2/cce/clusters/{clusterId}/apis/autoscaler.ccse.ctyun.cn/v1/horizontalnodeautoscaler/{name}  */
	RegionId string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&amp;api=5851&amp;data=87">资源池列表查询</a>  */
}

type CcseGetClusterAutoscalerPolicyResponse struct {
	StatusCode int32  `json:"statusCode,omitempty"` /*  状态码  */
	Message    string `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  string `json:"returnObj,omitempty"`  /*  返回结果  */
	Error      string `json:"error,omitempty"`      /*  错误码  */
}
