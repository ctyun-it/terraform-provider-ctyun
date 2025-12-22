package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcCgwBillNewApi
/* 云企业路由器按需订单订购 */
type EcEcCgwBillNewApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcCgwBillNewApi(client *core.CtyunClient) *EcEcCgwBillNewApi {
	return &EcEcCgwBillNewApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/cgw-bill/new",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcCgwBillNewApi) Do(ctx context.Context, credential core.Credential, req *EcEcCgwBillNewRequest) (*EcEcCgwBillNewResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcCgwBillNewRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EcEcCgwBillNewResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcCgwBillNewRequest struct {
	EcID            string  `json:"ecID"`                      /* 云间高速ID */
	RegionID        *string `json:"regionID,omitempty"`        /* 资源池ID */
	ClientToken     *string `json:"clientToken,omitempty"`     /* 客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一 */
	PayVoucherPrice *string `json:"payVoucherPrice,omitempty"` /* 代金券金额，只适用于预付费客户自动支付，若代金券支付金额传0或者控制符，则不适用代金券支付（小数会只保留2位，非负） */
}

type EcEcCgwBillNewResponse struct {
	StatusCode  *int32                           `json:"statusCode"`  /* 返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败 */
	ErrorCode   *string                          `json:"errorCode"`   /* 业务细分码，为product.module.code三段式码 */
	Message     *string                          `json:"message"`     /* 失败时的错误描述，一般为英文描述 */
	Description *string                          `json:"description"` /* 失败时的错误描述，一般为中文描述 */
	TraceID     *string                          `json:"traceID"`     /* 链路追踪ID */
	ReturnObj   *EcEcCgwBillNewReturnObjResponse `json:"returnObj"`   /* 返回参数 */
}

type EcEcCgwBillNewReturnObjResponse struct {
	MasterResourceStatus *string                          `json:"masterResourceStatus"` /* 主资源状态，只有主订单资源会返回 */
	RegionID             *string                          `json:"regionID"`             /* 资源所属资源池ID */
	MasterOrderID        *string                          `json:"masterOrderID"`        /* 订单ID。调用方在拿到masterOrderID之后，在若干错误情况下，可以使用materOrderID进一步确认订单状态及资源状态 */
	MasterResourceID     *string                          `json:"masterResourceID"`     /* 主资源ID */
	MasterOrderNO        *string                          `json:"masterOrderNO"`        /* 订单号 */
	Resources            []EcEcCgwBillNewResourceResponse `json:"resources"`            /* 资源明细列表 */
}

type EcEcCgwBillNewResourceResponse struct {
	OrderID          *string `json:"orderID"`          /* 订单ID */
	Status           *int32  `json:"status"`           /* 资源状态。参考masterResourceStatus */
	IsMaster         *bool   `json:"isMaster"`         /* 布尔类型，是否是主资源项 */
	ResourceType     *string `json:"resourceType"`     /* 本参数表示订单资源类型 */
	ResourceID       *string `json:"resourceID"`       /* 单项资源的变配、续订、退订等需要该资源项的ID */
	MasterOrderID    *string `json:"masterOrderID"`    /* 订单ID。调用方在拿到masterOrderID之后，在若干错误情况下，可以使用materOrderID进一步确认订单状态及资源状态 */
	UpdateTime       *int64  `json:"updateTime"`       /* 更新时刻，epoch时戳，毫秒精度 */
	ExpireTime       *int64  `json:"expireTime"`       /* 过期时刻，epoch时戳，毫秒精度 */
	MasterResourceID *string `json:"masterResourceID"` /* 主资源ID */
	ItemValue        *int32  `json:"itemValue"`        /* 资源规格，带宽包大小，单位MB */
	StartTime        *int64  `json:"startTime"`        /* 启动时刻，epoch时戳，毫秒精度 */
	CreateTime       *int64  `json:"createTime"`       /* 创建时刻，epoch时戳，毫秒精度 */
}
