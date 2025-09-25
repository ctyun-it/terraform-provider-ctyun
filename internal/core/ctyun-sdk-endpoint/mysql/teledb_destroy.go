package mysql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbDestroyApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbDestroyApi(client *ctyunsdk.CtyunClient) *TeledbDestroyApi {
	return &TeledbDestroyApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/teledb-acceptor/v2/openapi/accept-order-info/destroyOrder",
		},
	}
}

func (this *TeledbDestroyApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbDestroyRequest, header *TeledbDestroyRequestHeader) (destroyResp *TeledbDestroyResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != "" {
		builder.AddHeader("project-id", header.ProjectID)
	}
	if req.InstId == "" {
		err = errors.New("instId 为空")
		return
	}
	builder.AddParam("instId", req.InstId)
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	destroyResp = &TeledbDestroyResponse{}
	err = resp.Parse(destroyResp)
	if err != nil {
		return
	}
	return destroyResp, nil
}

type TeledbDestroyRequest struct {
	InstId string `json:"instId"` // 实例ID，必填
}
type TeledbDestroyRequestHeader struct {
	ProjectID string `json:"projectID"`
}
type TeledbDestroyResponse struct {
	StatusCode int32                           `json:"statusCode"`      // 接口状态码
	Error      *string                         `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                          `json:"message"`         // 描述信息
	ReturnObj  *TeledbDestroyResponseReturnObj `json:"returnObj"`       // 返回对象，类型为 DataObject
}

type TeledbDestroyResponseReturnObj struct {
	Data *TeledbDestroyResponseReturnObjData `json:"data"`
}
type TeledbDestroyResponseReturnObjData struct {
	ErrorMessage string  `json:"errorMessage"` // 错误内容
	Submitted    bool    `json:"submitted"`    // 是否成功
	NewOrderId   string  `json:"newOrderId"`   // 订单ID
	NewOrderNo   string  `json:"newOrderNo"`   // 订单号
	TotalPrice   float64 `json:"totalPrice"`   // 总价
}
