package pgsql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlGetIDByOrderApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlGetIDByOrderApi(client *ctyunsdk.CtyunClient) *PgsqlGetIDByOrderApi {
	return &PgsqlGetIDByOrderApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/teledb-acceptor/v2/openapi/accept-order-info/getInstIdByOrderId",
		},
	}
}

func (this *PgsqlGetIDByOrderApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlGetIDByOrderRequest, header *PgsqlGetIDByOrderRequestHeader) (destroyResp *PgsqlGetIDByOrderResponse, err error) {
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
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	destroyResp = &PgsqlGetIDByOrderResponse{}
	err = resp.Parse(destroyResp)
	if err != nil {
		return
	}
	return destroyResp, nil
}

type PgsqlGetIDByOrderRequest struct {
	OrderID string `json:"newOrderId"` // 实例ID，必填
}
type PgsqlGetIDByOrderRequestHeader struct {
	ProjectID string `json:"projectID"`
}
type PgsqlGetIDByOrderResponse struct {
	StatusCode int32                               `json:"statusCode"` // 接口状态码
	Error      *string                             `json:"error"`      // 错误码，失败时返回，成功时为空
	Message    string                              `json:"message"`    // 描述信息
	ReturnObj  *PgsqlGetIDByOrderResponseReturnObj `json:"returnObj"`  // 返回对象，类型为 DataObject
}

type PgsqlGetIDByOrderResponseReturnObj struct {
	Data []string `json:"data"`
}
