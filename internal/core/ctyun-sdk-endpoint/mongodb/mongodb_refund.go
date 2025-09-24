package mongodb

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type MongodbRefundApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewMongodbRefundApi(client *ctyunsdk.CtyunClient) *MongodbRefundApi {
	return &MongodbRefundApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/teledb-acceptor/v2/openapi/accept-order-info/refundOrder",
		},
	}
}

func (this *MongodbRefundApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *MongodbRefundRequest, header *MongodbRefundRequestHeader) (refundResp *MongodbRefundResponse, err error) {
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
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameMongodb, builder)
	if err != nil {
		return
	}
	refundResp = &MongodbRefundResponse{}
	err = resp.Parse(refundResp)
	if err != nil {
		return
	}
	return refundResp, nil
}

type MongodbRefundRequest struct {
	InstId string `json:"instId"` // 实例ID，必填
}
type MongodbRefundRequestHeader struct {
	ProjectID string `json:"projectID"`
}
type MongodbRefundResponse struct {
	StatusCode int32                           `json:"statusCode"`      // 接口状态码
	Error      *string                         `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                          `json:"message"`         // 描述信息
	ReturnObj  *MongodbRefundResponseReturnObj `json:"returnObj"`       // 返回对象，类型为 DataObject
}

type MongodbRefundResponseReturnObj struct {
	Data *MongodbRefundResponseReturnObjData `json:"data"`
}
type MongodbRefundResponseReturnObjData struct {
	ErrorMessage string  `json:"errorMessage"` // 错误内容
	Submitted    bool    `json:"submitted"`    // 是否成功
	NewOrderId   string  `json:"newOrderId"`   // 订单ID
	NewOrderNo   string  `json:"newOrderNo"`   // 订单号
	TotalPrice   float64 `json:"totalPrice"`   // 总价
}
