package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseBatchDeleteNodePoolApi
/* 调用该接口批量删除节点池。
 */type CcseBatchDeleteNodePoolApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseBatchDeleteNodePoolApi(client *core.CtyunClient) *CcseBatchDeleteNodePoolApi {
	return &CcseBatchDeleteNodePoolApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodDelete,
			UrlPath:      "/v2/cce/clusters/{clusterId}/nodepool/batchdelete",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseBatchDeleteNodePoolApi) Do(ctx context.Context, credential core.Credential, req *CcseBatchDeleteNodePoolRequest) (*CcseBatchDeleteNodePoolResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
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
	var resp CcseBatchDeleteNodePoolResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseBatchDeleteNodePoolRequest struct {
	ClusterId string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105" target="_blank">如何获取接口URI中参数</a>。  */
	RegionId  string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
	NodePoolNames []string `json:"nodePoolNames"` /*  节点池名称列表  */
}

type CcseBatchDeleteNodePoolResponse struct {
	StatusCode int32                                     `json:"statusCode,omitempty"` /*  状态码  */
	Message    string                                    `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  *CcseBatchDeleteNodePoolReturnObjResponse `json:"returnObj"`            /*  返回对象  */
	Error      string                                    `json:"error,omitempty"`      /*  错误码  */
}

type CcseBatchDeleteNodePoolReturnObjResponse struct {
	SuccessList  []string `json:"successList"`            /*  批量操作成功列表  */
	FailedList   []string `json:"failedList"`             /*  批量操作失败列表  */
	Total        int32    `json:"total,omitempty"`        /*  批量操作数据个数  */
	SuccessTotal int32    `json:"successTotal,omitempty"` /*  批量操作成功数据个数  */
	FailedTotal  int32    `json:"failedTotal,omitempty"`  /*  批量操作失败数据个数  */
}
