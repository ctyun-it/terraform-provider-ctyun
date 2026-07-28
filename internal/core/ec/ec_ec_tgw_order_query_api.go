package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcTgwOrderQueryApi
/* 按需订单查询 */
type EcEcTgwOrderQueryApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcTgwOrderQueryApi(client *core.CtyunClient) *EcEcTgwOrderQueryApi {
	return &EcEcTgwOrderQueryApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ec/tgw-order/query",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcTgwOrderQueryApi) Do(ctx context.Context, credential core.Credential, req *EcEcTgwOrderQueryRequest) (*EcEcTgwOrderQueryResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("ecID", req.EcID)
	if req.ResourceID != nil && *req.ResourceID != "" {
		ctReq.AddParam("resourceID", *req.ResourceID)
	}
	if req.OrderType != nil && *req.OrderType != "" {
		ctReq.AddParam("orderType", *req.OrderType)
	}
	if req.OrderState != nil && *req.OrderState != "" {
		ctReq.AddParam("orderState", *req.OrderState)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EcEcTgwOrderQueryResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcTgwOrderQueryRequest struct {
	EcID       string  `json:"ecID"`                 /* 云间高速ID */
	ResourceID *string `json:"resourceID,omitempty"` /* 资源ID */
	OrderType  *string `json:"orderType,omitempty"`  /* 订单类型 */
	OrderState *string `json:"orderState,omitempty"` /* 订单状态 */
}

type EcEcTgwOrderQueryResponse struct {
	StatusCode  *int32                            `json:"statusCode"`  /* 返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败 */
	ErrorCode   *string                           `json:"errorCode"`   /* 业务细分码，为product.module.code三段式码 */
	Message     *string                           `json:"message"`     /* 失败时的错误描述，一般为英文描述 */
	Description *string                           `json:"description"` /* 失败时的错误描述，一般为中文描述 */
	TraceID     *string                           `json:"traceID"`     /* 链路追踪ID */
	ReturnObj   *EcTgwOrderQueryReturnObjResponse `json:"returnObj"`   /* 返回参数 */
}

type EcTgwOrderQueryReturnObjResponse struct {
	CurrentCount *int32                            `json:"currentCount"` /* 当前页记录数 */
	TotalPage    *int32                            `json:"totalPage"`    /* 总页数 */
	TotalCount   *int32                            `json:"totalCount"`   /* 查询的总记录数 */
	Msg          *string                           `json:"msg"`          /* 消息 */
	Results      []*EcTgwOrderQueryResultsResponse `json:"results"`      /* 返回查询结果，Json数组 */
}

type EcTgwOrderQueryResultsResponse struct {
	OrderID    *string `json:"orderID"`    /* 订单ID */
	ExpireDate *string `json:"expireDate"` /* 过期时间 */
	ResourceID *string `json:"resourceID"` /* 资源ID */
	OrderType  *int32  `json:"orderType"`  /* 订单类型 */
	EcID       *string `json:"ecID"`       /* 云间高速ID */
	TenantID   *string `json:"tenantID"`   /* 用户ID */
	CreateDate *string `json:"createDate"` /* 创建时间 */
}
