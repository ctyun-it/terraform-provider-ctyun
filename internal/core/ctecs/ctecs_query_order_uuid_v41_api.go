package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsQueryOrderUuidV41Api
/* 根据订单号masterOrderId，输出订单状态、资源类型、资源uuid列表
 */type CtecsQueryOrderUuidV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsQueryOrderUuidV41Api(client *core.CtyunClient) *CtecsQueryOrderUuidV41Api {
	return &CtecsQueryOrderUuidV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/order/query-uuid",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsQueryOrderUuidV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsQueryOrderUuidV41Request) (*CtecsQueryOrderUuidV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("masterOrderId", req.MasterOrderId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsQueryOrderUuidV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsQueryOrderUuidV41Request struct {
	MasterOrderId string /*  订单id  */
}

type CtecsQueryOrderUuidV41Response struct {
	StatusCode  int32                                    `json:"statusCode,omitempty"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   string                                   `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码。为空表示成功。  */
	Message     string                                   `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                   `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsQueryOrderUuidV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
	Error       string                                   `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码。请求成功时不返回该字段  */
}

type CtecsQueryOrderUuidV41ReturnObjResponse struct {
	OrderStatus  string   `json:"orderStatus,omitempty"`  /*  订单状态  */
	ResourceType string   `json:"resourceType,omitempty"` /*  资源类型  */
	ResourceUuid []string `json:"resourceUuid"`           /*  资源uuid  */
}
