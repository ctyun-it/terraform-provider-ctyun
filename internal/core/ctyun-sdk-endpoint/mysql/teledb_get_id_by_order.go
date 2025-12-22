package mysql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbGetIDByOrderApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbGetIDByOrderApi(client *ctyunsdk.CtyunClient) *TeledbGetIDByOrderApi {
	return &TeledbGetIDByOrderApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/teledb-acceptor/v2/openapi/accept-order-info/getInstIdByOrderId",
		},
	}
}

func (this *TeledbGetIDByOrderApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbGetIDByOrderRequest, header *TeledbGetIDByOrderRequestHeader) (destroyResp *TeledbGetIDByOrderResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != "" {
		builder.AddHeader("project-id", header.ProjectID)
	}
	if req.OrderID == "" {
		err = errors.New("OrderID为空")
		return
	}
	builder.AddParam("newOrderId", req.OrderID)
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	destroyResp = &TeledbGetIDByOrderResponse{}
	err = resp.Parse(destroyResp)
	if err != nil {
		return
	}
	return destroyResp, nil
}

type TeledbGetIDByOrderRequest struct {
	OrderID string `json:"newOrderId"` // 实例ID，必填
}
type TeledbGetIDByOrderRequestHeader struct {
	ProjectID string `json:"projectID"`
}
type TeledbGetIDByOrderResponse struct {
	StatusCode int32                                `json:"statusCode"` // 接口状态码
	Error      *string                              `json:"error"`      // 错误码，失败时返回，成功时为空
	Message    string                               `json:"message"`    // 描述信息
	ReturnObj  *TeledbGetIDByOrderResponseReturnObj `json:"returnObj"`  // 返回对象，类型为 DataObject
}

type TeledbGetIDByOrderResponseReturnObj struct {
	Data []string `json:"data"`
}
