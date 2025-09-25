package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2TransChargeTypeApi
/* 变更分布式缓存Redis实例的付费类型按需转包周期。
 */type Dcs2TransChargeTypeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2TransChargeTypeApi(client *core.CtyunClient) *Dcs2TransChargeTypeApi {
	return &Dcs2TransChargeTypeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/res/spuInst/transChargeType",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2TransChargeTypeApi) Do(ctx context.Context, credential core.Credential, req *Dcs2TransChargeTypeRequest) (*Dcs2TransChargeTypeResponse, error) {
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
	var resp Dcs2TransChargeTypeResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2TransChargeTypeRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	SpuCode    string `json:"spuCode,omitempty"`    /*  产品类型  */
	CycleType  string `json:"cycleType,omitempty"`  /*  包周期类型<li>3：月<li>5：年  */
	CycleCnt   int32  `json:"cycleCnt,omitempty"`   /*  包周期数量，可选范围：1-6,12,24,36  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID  */
	AutoPay    bool   `json:"autoPay"`              /*  是否自动支付(仅对包周期实例有效)：<li>true：自动付费<li>false：手动付费(默认值)<br>说明：选择为手动付费时，您需要在控制台的右上角选择我的订单，然后单击左侧导航栏的订单管理 -> 待支付订单，找到目标订单进行支付。  */
}

type Dcs2TransChargeTypeResponse struct {
	StatusCode int32                                 `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2TransChargeTypeReturnObjResponse `json:"returnObj"`            /*  返回对象  */
	RequestId  string                                `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2TransChargeTypeReturnObjResponse struct {
	MasterOrderId   string `json:"masterOrderId,omitempty"`   /*  主订单ID  */
	MasterOrderNo   string `json:"masterOrderNo,omitempty"`   /*  主订单编码  */
	MasterOrderType string `json:"masterOrderType,omitempty"` /*  主订单类型  */
	AccountId       string `json:"accountId,omitempty"`       /*  账号ID  */
	UserId          string `json:"userId,omitempty"`          /*  用户ID  */
	Source          string `json:"source,omitempty"`          /*  资源编码  */
	CreateDate      string `json:"createDate,omitempty"`      /*  创建时间  */
	CreateStff      int32  `json:"createStff,omitempty"`      /*  创建标志  */
	UpdateStaff     int32  `json:"updateStaff,omitempty"`     /*  更新标志  */
	UpdateDate      string `json:"updateDate,omitempty"`      /*  更新时间  */
}
