package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
	"strings"
)

// EcsTagOnDemandApi 包周期付费云主机标记到期转按需
type EcsTagOnDemandApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEcsTagOnDemandApi(client *ctyunsdk.CtyunClient) *EcsTagOnDemandApi {
	return &EcsTagOnDemandApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ecs/tag-on-demand",
		},
	}
}

func (this *EcsTagOnDemandApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EcsTagOnDemandRequest) (*EcsTagOnDemandResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&ecsTagOnDemandRealRequest{
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

	var realResponse ecsTagOnDemandOrderInfoRealResponse
	err = response.ParseByStandardModelWithCheck(&realResponse)
	if err != nil {
		return nil, err
	}
	var orderInfo []EcsTagOnDemandOrderInfoResponse
	for _, order := range realResponse.OrderInfo {
		orderInfo = append(orderInfo, EcsTagOnDemandOrderInfoResponse{
			OrderId: order.OrderID,
			OrderNo: order.OrderNo,
		})
	}
	return &EcsTagOnDemandResponse{
		OrderInfo: orderInfo,
	}, nil
}

type ecsTagOnDemandRealRequest struct {
	ClientToken string `json:"clientToken"`
	RegionID    string `json:"regionID"`
	InstanceIDs string `json:"instanceIDs"`
}

type ecsTagOnDemandOrderInfoRealResponse struct {
	OrderInfo []struct {
		OrderID string `json:"orderID"`
		OrderNo string `json:"orderNo"`
	} `json:"orderInfo"`
}

type EcsTagOnDemandRequest struct {
	ClientToken string
	RegionId    string
	InstanceIds []string
}

type EcsTagOnDemandOrderInfoResponse struct {
	OrderId string
	OrderNo string
}

type EcsTagOnDemandResponse struct {
	OrderInfo []EcsTagOnDemandOrderInfoResponse
}
