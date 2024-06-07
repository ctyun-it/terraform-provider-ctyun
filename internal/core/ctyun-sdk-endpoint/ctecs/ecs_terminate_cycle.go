package ctecs

import (
	"context"
	"net/http"
	"strings"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// EcsTerminateCycleApi 包周期终止
type EcsTerminateCycleApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEcsTerminateCycleApi(client *ctyunsdk.CtyunClient) *EcsTerminateCycleApi {
	return &EcsTerminateCycleApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ecs/terminate-cycle",
		},
	}
}

func (this *EcsTerminateCycleApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EcsTerminateCycleRequest) (*EcsTerminateCycleResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&ecsTerminateCycleRealRequest{
		ClientToken: req.ClientToken,
		RegionID:    req.RegionId,
		InstanceIDs: strings.Join(req.InstanceIds, ","),
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtecs, builder)
	if err != nil {
		return nil, err
	}

	var realResponse ecsTerminateCycleOrderInfoRealResponse
	err = response.ParseByStandardModelWithCheck(&realResponse)
	if err != nil {
		return nil, err
	}
	var orderInfo []EcsTerminateCycleOrderInfoResponse
	for _, order := range realResponse.OrderInfo {
		orderInfo = append(orderInfo, EcsTerminateCycleOrderInfoResponse{
			OrderId: order.OrderID,
			OrderNo: order.OrderNo,
		})
	}
	return &EcsTerminateCycleResponse{
		OrderInfo: orderInfo,
	}, nil
}

type ecsTerminateCycleRealRequest struct {
	ClientToken string `json:"clientToken"`
	RegionID    string `json:"regionID"`
	InstanceIDs string `json:"instanceIDs"`
}

type ecsTerminateCycleOrderInfoRealResponse struct {
	OrderInfo []struct {
		OrderID string `json:"orderID"`
		OrderNo string `json:"orderNo"`
	} `json:"orderInfo"`
}

type EcsTerminateCycleRequest struct {
	ClientToken string
	RegionId    string
	InstanceIds []string
}

type EcsTerminateCycleOrderInfoResponse struct {
	OrderId string
	OrderNo string
}

type EcsTerminateCycleResponse struct {
	OrderInfo []EcsTerminateCycleOrderInfoResponse
}
