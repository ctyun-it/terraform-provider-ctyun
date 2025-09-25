package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2QueryInstAutoRenewStatusApi
/* 查看实例是否开通自动续费。
 */type Dcs2QueryInstAutoRenewStatusApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2QueryInstAutoRenewStatusApi(client *core.CtyunClient) *Dcs2QueryInstAutoRenewStatusApi {
	return &Dcs2QueryInstAutoRenewStatusApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/res/spuInst/queryInstAutoRenewStatus",
			ContentType:  "",
		},
	}
}

func (a *Dcs2QueryInstAutoRenewStatusApi) Do(ctx context.Context, credential core.Credential, req *Dcs2QueryInstAutoRenewStatusRequest) (*Dcs2QueryInstAutoRenewStatusResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("spuCode", req.SpuCode)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2QueryInstAutoRenewStatusResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2QueryInstAutoRenewStatusRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例ID  */
	SpuCode    string /*  产品类型，DCS2-缓存  */
}

type Dcs2QueryInstAutoRenewStatusResponse struct {
	StatusCode int32                                          `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                         `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2QueryInstAutoRenewStatusReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                         `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                         `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                         `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2QueryInstAutoRenewStatusReturnObjResponse struct {
	ErrorMsg      string                                                        `json:"errorMsg,omitempty"` /*  错误信息  */
	ResourcesDTOs []*Dcs2QueryInstAutoRenewStatusReturnObjResourcesDTOsResponse `json:"resourcesDTOs"`      /*  数组  */
}

type Dcs2QueryInstAutoRenewStatusReturnObjResourcesDTOsResponse struct {
	InstanceId       string `json:"instanceId,omitempty"`       /*  虚拟实例标识  */
	MasterResourceId string `json:"masterResourceId,omitempty"` /*  主虚拟资源ID  */
	ResourceType     string `json:"resourceType,omitempty"`     /*  资源类型  */
	StartDate        string `json:"startDate,omitempty"`        /*  开始时间,值为 Unix 时间戳（毫秒级）  */
	ExpireDate       string `json:"expireDate,omitempty"`       /*  到期时间,值为 Unix 时间戳（毫秒级）  */
	AutoRenewStatus  string `json:"autoRenewStatus,omitempty"`  /*  自动续订状态<li>0/null：不自动续订<li>1：自动续订  */
}
