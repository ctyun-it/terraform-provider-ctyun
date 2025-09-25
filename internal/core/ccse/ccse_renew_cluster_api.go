package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseRenewClusterApi
/* 续订集群
 */type CcseRenewClusterApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseRenewClusterApi(client *core.CtyunClient) *CcseRenewClusterApi {
	return &CcseRenewClusterApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/cce/clusters/renew",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseRenewClusterApi) Do(ctx context.Context, credential core.Credential, req *CcseRenewClusterRequest) (*CcseRenewClusterResponse, error) {
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
	var resp CcseRenewClusterResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseRenewClusterRequest struct {
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
	CycleType  string `json:"cycleType,omitempty"`  /*  订购周期类型，3（按月）、5（按年）  */
	CycleCnt   int32  `json:"cycleCnt,omitempty"`   /*  订购时长  */
}

type CcseRenewClusterResponse struct {
	StatusCode int32                              `json:"statusCode,omitempty"` /*  响应状态码  */
	Message    string                             `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *CcseRenewClusterReturnObjResponse `json:"returnObj"`            /*  响应对象  */
	Error      string                             `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type CcseRenewClusterReturnObjResponse struct {
	OrderNo string `json:"orderNo,omitempty"` /*  20240925163106798886  */
	OrderId string `json:"orderId,omitempty"` /*  8a0109f9196d41db9d92893fe4522c75  */
}
