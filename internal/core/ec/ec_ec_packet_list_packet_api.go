package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcPacketListPacketApi
/* 查询带宽包 */
type EcEcPacketListPacketApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcPacketListPacketApi(client *core.CtyunClient) *EcEcPacketListPacketApi {
	return &EcEcPacketListPacketApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ec/bandwidth-packet/list",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcPacketListPacketApi) Do(ctx context.Context, credential core.Credential, req *EcEcPacketListPacketRequest) (*EcEcPacketListPacketResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("ecID", req.EcID)
	if req.PacketID != nil && *req.PacketID != "" {
		ctReq.AddParam("packetID", *req.PacketID)
	}
	if req.ResourceID != nil && *req.ResourceID != "" {
		ctReq.AddParam("resourceID", *req.ResourceID)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EcEcPacketListPacketResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcPacketListPacketRequest struct {
	EcID       string  `json:"ecID"`                 /*  云间高速实例ID  */
	PacketID   *string `json:"packetID,omitempty"`   /*  带宽包ID  */
	ResourceID *string `json:"resourceID,omitempty"` /*  带宽包 resource ID  */
}

type EcEcPacketListPacketResponse struct {
	StatusCode  *int32                               `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                              `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                              `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                              `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                              `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcPacketListPacketReturnObjResponse `json:"returnObj"`   /*  返回参数  */
}

type EcPacketListPacketReturnObjResponse struct {
	CurrentCount *int32                               `json:"currentCount"` /* 当前页记录数 */
	TotalPage    *int32                               `json:"totalPage"`    /* 总页数 */
	TotalCount   *int32                               `json:"totalCount"`   /* 查询的总记录数 */
	Results      []*EcPacketListPacketResultsResponse `json:"results"`      /* 返回查询结果，Json数组 */
}

type EcPacketListPacketResultsResponse struct {
	PacketID     *string `json:"packetID"`     /* 带宽包ID */
	PacketName   *string `json:"packetName"`   /* 带宽包名称 */
	ResourceID   *string `json:"resourceID"`   /* 订单资源ID */
	Rate         *int32  `json:"rate"`         /* 带宽（MB） */
	Status       *string `json:"status"`       /* 运行状态 */
	AreaA        *string `json:"areaA"`        /* 区域A类型 */
	AreaB        *string `json:"areaB"`        /* 区域B类型 */
	EcID         *string `json:"ecID"`         /* 云间高速实例ID */
	CreateDate   *string `json:"createDate"`   /* 创建时间 */
	DeleteDate   *string `json:"deleteDate"`   /* 到期时间 */
	BillingModel *string `json:"billingModel"` /* 计费类型 */
	UsedRate     *int32  `json:"usedRate"`     /* 已用带宽（MB） */
	UsableRate   *int32  `json:"usableRate"`   /* 剩余可用带宽（MB） */
	ResourceType *int32  `json:"resourceType"` /* 带宽包类型 */
}
