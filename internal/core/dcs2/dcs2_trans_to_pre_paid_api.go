package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2TransToPrePaidApi
/* 分布式缓存Redis实例包周期转按需付费，到期后自动转为按需付费。
 */type Dcs2TransToPrePaidApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2TransToPrePaidApi(client *core.CtyunClient) *Dcs2TransToPrePaidApi {
	return &Dcs2TransToPrePaidApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/res/spuInst/transToPrePaid",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2TransToPrePaidApi) Do(ctx context.Context, credential core.Credential, req *Dcs2TransToPrePaidRequest) (*Dcs2TransToPrePaidResponse, error) {
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
	var resp Dcs2TransToPrePaidResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2TransToPrePaidRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID  */
	SpuCode    string `json:"spuCode,omitempty"`    /*  产品编码  */
}

type Dcs2TransToPrePaidResponse struct {
	StatusCode int32                                `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                               `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2TransToPrePaidReturnObjResponse `json:"returnObj"`            /*  响应对象  */
	RequestId  string                               `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                               `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                               `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2TransToPrePaidReturnObjResponse struct {
	MasterOrderId   string `json:"masterOrderId,omitempty"`   /*  主订单ID  */
	MasterOrderNo   string `json:"masterOrderNo,omitempty"`   /*  主订单编码  */
	MasterOrderType string `json:"masterOrderType,omitempty"` /*  主订单类型  */
	AccountId       string `json:"accountId,omitempty"`       /*  账号ID  */
	UserId          string `json:"userId,omitempty"`          /*  用户ID  */
	Source          string `json:"source,omitempty"`          /*  资源编码  */
	CreateDate      string `json:"createDate,omitempty"`      /*  创建时间  */
	UpdateDate      string `json:"updateDate,omitempty"`      /*  更新时间  */
}
