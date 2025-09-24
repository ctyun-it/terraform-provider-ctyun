package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// EcsDestroyInstanceApi
type EcsDestroyInstanceApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEcsDestroyInstanceApi(client *ctyunsdk.CtyunClient) *EcsDestroyInstanceApi {
	return &EcsDestroyInstanceApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ecs/destroy-instance",
		},
	}
}

func (this *EcsDestroyInstanceApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EcsDestroyInstanceRequest) (*EcsDestroyInstanceResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)

	_, err := builder.WriteJson(&EcsDestroyInstanceRealRequest{
		ClientToken: req.ClientToken,
		RegionID:    req.RegionID,
		InstanceID:  req.InstanceID,
	})

	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtecs, builder)
	if err != nil {
		return nil, err
	}

	var realResponse EcsDestroyInstanceRealResponse
	err = response.ParseByStandardModelWithCheck(&realResponse)
	if err != nil {
		return nil, err
	}

	return &EcsDestroyInstanceResponse{
		MasterOrderID: realResponse.MasterOrderID,
		MasterOrderNO: realResponse.MasterOrderNO,
		RegionID:      realResponse.RegionID,
	}, nil
}

type EcsDestroyInstanceRealRequest struct {
	ClientToken string `json:"clientToken,omitempty"`
	RegionID    string `json:"regionID,omitempty"`
	InstanceID  string `json:"instanceID,omitempty"`
}

type EcsDestroyInstanceRequest struct {
	ClientToken string
	RegionID    string
	InstanceID  string
}

type EcsDestroyInstanceRealResponse struct {
	MasterOrderID string `json:"masterOrderID,omitempty"`
	MasterOrderNO string `json:"masterOrderNO,omitempty"`
	RegionID      string `json:"regionID,omitempty"`
}

type EcsDestroyInstanceResponse struct {
	MasterOrderID string
	MasterOrderNO string
	RegionID      string
}
