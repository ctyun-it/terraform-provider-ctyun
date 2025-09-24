package pgsql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlRefundApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlRefundApi(client *ctyunsdk.CtyunClient) *PgsqlRefundApi {
	return &PgsqlRefundApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/teledb-acceptor/v2/openapi/accept-order-info/refundOrder",
		},
	}
}

func (this *PgsqlRefundApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlRefundRequest, header *PgsqlRefundRequestHeader) (refundResponse *PgsqlRefundResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}
	if req.InstId == "" {
		err = errors.New("missing required field: InstId")
		return
	}
	builder.AddParam("instId", req.InstId)
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	response := PgsqlRefundResponse{}
	err = resp.Parse(&response)
	if err != nil {
		return
	}
	return &response, nil
}

type PgsqlRefundRequest struct {
	InstId string `json:"instId"` //数据库实例ID
}

type PgsqlRefundRequestHeader struct {
	ProjectID *string `json:"projectId,omitempty"` //项目id
}

type PgsqlRefundResponse struct {
	StatusCode int32                         `json:"statusCode"` // 接口状态码，参考下方状态码
	Error      string                        `json:"error"`      // 错误码
	Message    string                        `json:"message"`    // 描述信息
	ReturnObj  *PgsqlRefundResponseReturnObj `json:"returnObj"`  // 返回对象，包含具体的返回数据
}

type PgsqlRefundResponseReturnObj struct {
	ErrorMessage string  `json:"errorMessage"` // 错误内容
	Submitted    bool    `json:"submitted"`    // 是否成功
	NewOrderId   string  `json:"newOrderId"`   // 订单ID
	NewOrderNo   string  `json:"newOrderNo"`   // 订单号
	TotalPrice   float64 `json:"totalPrice"`   // 总价
}
