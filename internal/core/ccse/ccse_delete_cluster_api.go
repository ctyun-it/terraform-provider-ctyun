package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseDeleteClusterApi
/* 调用该接口退订Kubernetes集群。
 */type CcseDeleteClusterApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseDeleteClusterApi(client *core.CtyunClient) *CcseDeleteClusterApi {
	return &CcseDeleteClusterApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/cce/clusters/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseDeleteClusterApi) Do(ctx context.Context, credential core.Credential, req *CcseDeleteClusterRequest) (*CcseDeleteClusterResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
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
	var resp CcseDeleteClusterResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseDeleteClusterRequest struct {
	RegionId string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>
	*/
	ResPoolId string `json:"resPoolId,omitempty"` /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>
	*/
	ProdInstId string `json:"prodInstId,omitempty"` /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
}

type CcseDeleteClusterResponse struct {
	StatusCode int32                               `json:"statusCode,omitempty"` /*  状态码  */
	Message    string                              `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  *CcseDeleteClusterReturnObjResponse `json:"returnObj"`            /*  返回对象  */
	Error      string                              `json:"error,omitempty"`      /*  返回错误码  */
	RequestId  string                              `json:"requestId"`
}

type CcseDeleteClusterReturnObjResponse struct {
	OrderId string `json:"orderId,omitempty"` /*  订单ID  */
	OrderNo string `json:"orderNo,omitempty"` /*  订单编码  */
}
