package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// EcsOrderQueryUuidApi 根据订单号查询uuid
// https://www.ctyun.cn/document/10026730/10069118
type EcsOrderQueryUuidApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEcsOrderQueryUuid(client *ctyunsdk.CtyunClient) *EcsOrderQueryUuidApi {
	return &EcsOrderQueryUuidApi{
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v4/ecs/order/query-uuid",
		},
		client: client,
	}
}

func (this *EcsOrderQueryUuidApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EcsOrderQueryUuidRequest) (*EcsOrderQueryUuidResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	builder.AddParam("masterOrderID", req.MasterOrderId)
	resp, requestError := this.client.RequestToEndpoint(ctx, EndpointNameCtecs, builder)
	if requestError != nil {
		return nil, requestError
	}
	var realResponse EcsOrderQueryUuidRealResponse
	err := resp.ParseByStandardModelWithCheck(&realResponse)
	if err != nil {
		return nil, err
	}
	return &EcsOrderQueryUuidResponse{
		InstanceIDList: realResponse.InstanceIDList,
		OrderStatus:    realResponse.OrderStatus,
	}, nil
}

type EcsOrderQueryUuidRealResponse struct {
	InstanceIDList []string `json:"instanceIDList"`
	OrderStatus    string   `json:"orderStatus"`
}

type EcsOrderQueryUuidRequest struct {
	MasterOrderId string
}

type EcsOrderQueryUuidResponse struct {
	InstanceIDList []string `json:"instanceIDList"`
	OrderStatus    string   `json:"orderStatus"`
}
