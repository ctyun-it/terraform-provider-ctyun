package ctecs

import (
	"context"
	"net/http"
	"strings"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// EcsChangeToCycleApi 按需付费转包周期
type EcsChangeToCycleApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEcsChangeToCycleApi(client *ctyunsdk.CtyunClient) *EcsChangeToCycleApi {
	return &EcsChangeToCycleApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ecs/change-to-cycle",
		},
	}
}

func (this *EcsChangeToCycleApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EcsChangeToCycleRequest) (*EcsChangeToCycleResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&ecsChangeToCycleRealRequest{
		ClientToken: req.ClientToken,
		RegionID:    req.RegionId,
		InstanceIDs: strings.Join(req.InstanceIds, ","),
		CycleType:   req.CycleType,
		CycleCount:  req.CycleCount,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtecs, builder)
	if err != nil {
		return nil, err
	}

	var realResponse ecsChangeToCycleRealResponse
	err = response.ParseByStandardModelWithCheck(&realResponse)
	if err != nil {
		return nil, err
	}
	var orderInfo []EcsChangeToCycleOrderInfoResponse
	for _, order := range realResponse.OrderInfo {
		orderInfo = append(orderInfo, EcsChangeToCycleOrderInfoResponse{
			OrderId: order.OrderID,
			OrderNo: order.OrderNo,
		})
	}
	return &EcsChangeToCycleResponse{
		OrderInfo: orderInfo,
	}, nil
}

type ecsChangeToCycleRealRequest struct {
	ClientToken string `json:"clientToken"`
	RegionID    string `json:"regionID"`
	InstanceIDs string `json:"instanceIDs"`
	CycleType   string `json:"cycleType"`
	CycleCount  int    `json:"cycleCount"`
}

type ecsChangeToCycleRealResponse struct {
	OrderInfo []struct {
		OrderID string `json:"orderID"`
		OrderNo string `json:"orderNo"`
	} `json:"orderInfo"`
}

type EcsChangeToCycleRequest struct {
	ClientToken string
	RegionId    string
	InstanceIds []string
	CycleType   string
	CycleCount  int
}

type EcsChangeToCycleOrderInfoResponse struct {
	OrderId string
	OrderNo string
}

type EcsChangeToCycleResponse struct {
	OrderInfo []EcsChangeToCycleOrderInfoResponse
}
