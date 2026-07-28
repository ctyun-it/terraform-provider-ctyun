package mongodb

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type MongodbGetIDByOrderApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewMongodbGetIDByOrderApi(client *ctyunsdk.CtyunClient) *MongodbGetIDByOrderApi {
	return &MongodbGetIDByOrderApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/teledb-acceptor/v2/openapi/accept-order-info/getInstIdByOrderId",
		},
	}
}

func (this *MongodbGetIDByOrderApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *MongodbGetIDByOrderRequest, header *MongodbGetIDByOrderRequestHeader) (destroyResp *MongodbGetIDByOrderResponse, err error) {
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
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameMongodb, builder)
	if err != nil {
		return
	}
	destroyResp = &MongodbGetIDByOrderResponse{}
	err = resp.Parse(destroyResp)
	if err != nil {
		return
	}
	return destroyResp, nil
}

type MongodbGetIDByOrderRequest struct {
	OrderID string `json:"newOrderId"` // 实例ID，必填
}
type MongodbGetIDByOrderRequestHeader struct {
	ProjectID string `json:"projectID"`
}
type MongodbGetIDByOrderResponse struct {
	StatusCode int32                                 `json:"statusCode"` // 接口状态码
	Error      *string                               `json:"error"`      // 错误码，失败时返回，成功时为空
	Message    string                                `json:"message"`    // 描述信息
	ReturnObj  *MongodbGetIDByOrderResponseReturnObj `json:"returnObj"`  // 返回对象，类型为 DataObject
}

type MongodbGetIDByOrderResponseReturnObj struct {
	Data []string `json:"data"`
}
