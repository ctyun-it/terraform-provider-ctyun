package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsQueryInstancIDV41Api
/* 根据输入订单号masterOrderID查询出云主机的UUID<br/><b>准备工作：</b><br/>&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权<br /><b>注意事项：</b><br />&emsp;&emsp;订单查询：通过创建云主机接口得到对应的主订单ID后，使用该接口获取订单状态，并在订单状态为成功之后获取对应的云主机ID
 */type CtecsQueryInstancIDV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsQueryInstancIDV41Api(client *core.CtyunClient) *CtecsQueryInstancIDV41Api {
	return &CtecsQueryInstancIDV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ecs/order/query-uuid",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsQueryInstancIDV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsQueryInstancIDV41Request) (*CtecsQueryInstancIDV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("masterOrderID", req.MasterOrderID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsQueryInstancIDV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsQueryInstancIDV41Request struct {
	MasterOrderID string /*  订单ID  */
}

type CtecsQueryInstancIDV41Response struct {
	StatusCode  int32                                    `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                   `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                   `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                   `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                   `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsQueryInstancIDV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsQueryInstancIDV41ReturnObjResponse struct {
	OrderStatus    string   `json:"orderStatus,omitempty"` /*  订单状态，详见枚举值表格  */
	InstanceIDList []string `json:"instanceIDList"`        /*  云主机的ID列表，订单处于创建中，返回为空列表。待订单完成后才能正常返回ID  */
}
