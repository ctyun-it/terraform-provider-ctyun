package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2DeleteInstancesApi
/* 批量退订多个分布式缓存Redis实例。
 */type Dcs2DeleteInstancesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DeleteInstancesApi(client *core.CtyunClient) *Dcs2DeleteInstancesApi {
	return &Dcs2DeleteInstancesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/lifeCycleServant/deleteInstances",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2DeleteInstancesApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DeleteInstancesRequest) (*Dcs2DeleteInstancesResponse, error) {
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
	var resp Dcs2DeleteInstancesResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DeleteInstancesRequest struct {
	RegionId    string   /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstIds []string `json:"prodInstIds"` /*  实例Id列表，最多允许10个  */
}

type Dcs2DeleteInstancesResponse struct {
	StatusCode int32                                 `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2DeleteInstancesReturnObjResponse `json:"returnObj"`            /*  响应对象  */
	RequestId  string                                `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DeleteInstancesReturnObjResponse struct {
	ErrorMessage               string                                                            `json:"errorMessage,omitempty"`     /*  错误信息  */
	BatchOrderPlacementResults []*Dcs2DeleteInstancesReturnObjBatchOrderPlacementResultsResponse `json:"batchOrderPlacementResults"` /*  OrderPlacementResult  */
}

type Dcs2DeleteInstancesReturnObjBatchOrderPlacementResultsResponse struct {
	ErrorMessage      string                                                                             `json:"errorMessage,omitempty"` /*  错误信息  */
	Submitted         *bool                                                                              `json:"submitted"`              /*  是否已提交  */
	OrderPlacedEvents []*Dcs2DeleteInstancesReturnObjBatchOrderPlacementResultsOrderPlacedEventsResponse `json:"orderPlacedEvents"`      /*  orderPlacedEvent  */
}

type Dcs2DeleteInstancesReturnObjBatchOrderPlacementResultsOrderPlacedEventsResponse struct {
	ErrorMessage string  `json:"errorMessage,omitempty"` /*  错误信息  */
	Submitted    *bool   `json:"submitted"`              /*  是否已提交  */
	NewOrderId   string  `json:"newOrderId,omitempty"`   /*  新订单ID  */
	NewOrderNo   string  `json:"newOrderNo,omitempty"`   /*  新订单号  */
	TotalPrice   float64 `json:"totalPrice"`             /*  总价  */
}
