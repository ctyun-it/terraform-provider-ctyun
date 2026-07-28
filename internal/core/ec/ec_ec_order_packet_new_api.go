package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcOrderPacketNewApi
/* 支持按需包年/包月订购云间高速带宽包 */
type EcEcOrderPacketNewApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcOrderPacketNewApi(client *core.CtyunClient) *EcEcOrderPacketNewApi {
	return &EcEcOrderPacketNewApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/packet/new",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcOrderPacketNewApi) Do(ctx context.Context, credential core.Credential, req *EcEcOrderPacketNewRequest) (*EcEcOrderPacketNewResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcOrderPacketNewRequest
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
	var resp EcEcOrderPacketNewResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcOrderPacketNewRequest struct {
	ClientToken     *string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一  */
	RegionID        string  `json:"regionID"`              /*  资源池ID  */
	PacketName      string  `json:"packetName"`            /*  带宽包名字  */
	EcID            string  `json:"ecID"`                  /*  云间高速ID  */
	Bandwidth       int32   `json:"bandwidth"`             /*  带宽，单位MB  */
	OnDemand        bool    `json:"onDemand"`              /*  布尔类型，是否按需下单。默认为false  */
	CycleType       string  `json:"cycleType"`             /*  包周期类型,当onDemand为False时，必须指定<br/>取值如下：<br/>'YEAR': 包年<br/>'MONTH':包月  */
	CycleCount      int32   `json:"cycleCount"`            /*  包周期数。onDemand为False时必须指定。周期最大长度不能超过36个月  */
	AreaA           *string `json:"areaA"`                 /*  区域A类型，取值如下：<br/>"china": 中国大陆,<br/>"APAC":亚太<br/>默认china   */
	AreaB           *string `json:"areaB"`                 /*  区域B类型，取值如下：<br/>"china": 中国大陆,<br/>"APAC":亚太<br/>默认china  */
	PayVoucherPrice *string `json:"payVoucherPrice"`       /*  代金券金额，只适用于预付费客户自动支付，若代金券支付金额传0或者控制符，则不适用代金券支付（小数会只保留2位，非负）  */
}

type EcEcOrderPacketNewResponse struct {
	StatusCode  *int32                               `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	Message     *string                              `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                              `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *EcEcOrderPacketNewReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   *string                              `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Details     *EcEcOrderPacketNewDetailsResponse   `json:"details"`     /*  错误明细。一般情况下，会对订单侧(bss)的订单业务相关的错误做明确的错误映射和提升，有唯一对应的errorCode  */
	Error       *string                              `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type EcEcOrderPacketNewReturnObjResponse struct {
	MasterOrderID        *string                                         `json:"masterOrderID"`        /*  订单ID。调用方在拿到masterOrderID之后，在若干错误情况下，可以使用materOrderID进一步确认订单状态及资源状态  */
	MasterOrderNO        *string                                         `json:"masterOrderNO"`        /*  订单号  */
	MasterResourceID     *string                                         `json:"masterResourceID"`     /*  主资源ID  */
	MasterResourceStatus *string                                         `json:"masterResourceStatus"` /*  主资源状态，只有主订单资源会返回  */
	RegionID             *string                                         `json:"regionID"`             /*  资源所属资源池ID  */
	Resources            []*EcEcOrderPacketNewReturnObjResourcesResponse `json:"resources"`            /*  资源明细列表,参考表resources  */
}

type EcEcOrderPacketNewDetailsResponse struct{}

type EcEcOrderPacketNewReturnObjResourcesResponse struct {
	ResourceID       *string `json:"resourceID"`       /*  单项资源的变配、续订、退订等需要该资源项的ID  */
	OrderID          *string `json:"orderID"`          /*  订单ID  */
	StartTime        *int64  `json:"startTime"`        /*  启动时刻，epoch时戳，毫秒精度  */
	ExpireTime       *int64  `json:"expireTime"`       /*  过期时刻，epoch时戳，毫秒精度  */
	CreateTime       *int64  `json:"createTime"`       /*  创建时刻，epoch时戳，毫秒精度  */
	UpdateTime       *int64  `json:"updateTime"`       /*  更新时刻，epoch时戳，毫秒精度  */
	Status           *int32  `json:"status"`           /*  资源状态。参考masterResourceStatus  */
	IsMaster         *bool   `json:"isMaster"`         /*  布尔类型，是否是主资源项  */
	ItemValue        *int32  `json:"itemValue"`        /*  资源规格，带宽包大小，单位MB  */
	ResourceType     *string `json:"resourceType"`     /*  本参数表示订单资源类型  */
	MasterOrderID    *string `json:"masterOrderID"`    /*  订单ID。调用方在拿到masterOrderID之后，在若干错误情况下，可以使用materOrderID进一步确认订单状态及资源状态  */
	MasterResourceID *string `json:"masterResourceID"` /*  主资源ID  */
}
