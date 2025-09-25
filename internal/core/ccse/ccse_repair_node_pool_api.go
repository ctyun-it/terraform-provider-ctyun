package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseRepairNodePoolApi
/* 调用该接口创建节点池修复运维任务。
 */type CcseRepairNodePoolApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseRepairNodePoolApi(client *core.CtyunClient) *CcseRepairNodePoolApi {
	return &CcseRepairNodePoolApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPut,
			UrlPath:      "/v2/cce/clusters/{clusterId}/nodepool/{nodePoolId}/repair",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseRepairNodePoolApi) Do(ctx context.Context, credential core.Credential, req *CcseRepairNodePoolRequest) (*CcseRepairNodePoolResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder = builder.ReplaceUrl("nodePoolId", req.NodePoolId)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseRepairNodePoolResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseRepairNodePoolRequest struct {
	ClusterId  string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105" target="_blank">如何获取接口URI中参数</a>。  */
	NodePoolId string /*  节点池ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105" target="_blank">如何获取接口URI中参数</a>。  */
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
	NodeNames []string `json:"nodeNames"` /*  需要修复节点名称  */
}

type CcseRepairNodePoolResponse struct {
	StatusCode int32                                `json:"statusCode,omitempty"` /*  响应状态码  */
	Message    string                               `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *CcseRepairNodePoolReturnObjResponse `json:"returnObj"`            /*  响应对象  */
	Error      string                               `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type CcseRepairNodePoolReturnObjResponse struct {
	TaskId string `json:"taskId,omitempty"` /*  任务ID  */
}
