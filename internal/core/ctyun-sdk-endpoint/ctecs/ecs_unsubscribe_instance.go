package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// EcsUnsubscribeInstanceApi 释放云主机
// https://www.ctyun.cn/document/10026730/10106596
type EcsUnsubscribeInstanceApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEcsUnsubscribeInstanceApi(client *ctyunsdk.CtyunClient) *EcsUnsubscribeInstanceApi {
	return &EcsUnsubscribeInstanceApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ecs/unsubscribe-instance",
		},
	}
}

func (this *EcsUnsubscribeInstanceApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EcsUnsubscribeInstanceRequest) (*EcsUnsubscribeInstanceResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&ecsUnsubscribeInstanceRealRequest{
		ClientToken: req.ClientToken,
		RegionID:    req.RegionId,
		InstanceID:  req.InstanceId,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtecs, builder)
	if err != nil {
		return nil, err
	}

	var realResponse ecsUnsubscribeInstanceRealResponse
	err = response.ParseByStandardModelWithCheck(&realResponse)
	if err != nil {
		return nil, err
	}
	return &EcsUnsubscribeInstanceResponse{
		MasterOrderNo: realResponse.MasterOrderNO,
		RegionId:      realResponse.RegionID,
		MasterOrderId: realResponse.MasterOrderID,
	}, nil
}

type ecsUnsubscribeInstanceRealRequest struct {
	ClientToken string `json:"clientToken,omitempty"`
	RegionID    string `json:"regionID,omitempty"`
	InstanceID  string `json:"instanceID,omitempty"`
}

type EcsUnsubscribeInstanceRequest struct {
	ClientToken string
	RegionId    string
	InstanceId  string
}

type ecsUnsubscribeInstanceRealResponse struct {
	MasterOrderNO string `json:"masterOrderNO"`
	RegionID      string `json:"regionID"`
	MasterOrderID string `json:"masterOrderID"`
}

type EcsUnsubscribeInstanceResponse struct {
	MasterOrderNo string
	RegionId      string
	MasterOrderId string
}
