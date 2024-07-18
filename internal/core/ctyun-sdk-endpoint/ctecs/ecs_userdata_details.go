package ctecs

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// EcsUserdataDetailsApi
type EcsUserdataDetailsApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEcsUserdataDetailsApi(client *ctyunsdk.CtyunClient) *EcsUserdataDetailsApi {
	return &EcsUserdataDetailsApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v4/ecs/userdata/details",
		},
	}
}

func (this *EcsUserdataDetailsApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EcsUserdataDetailsRequest) (*EcsUserdataDetailsResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	builder.AddParam("regionID", req.RegionID)
	builder.AddParam("instanceID", req.InstanceID)

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtecs, builder)
	if err != nil {
		return nil, err
	}

	var realResponse EcsUserdataDetailsRealResponse
	err = response.ParseByStandardModelWithCheck(&realResponse)
	if err != nil {
		return nil, err
	}

	return &EcsUserdataDetailsResponse{
		Userdata: realResponse.Userdata,
	}, nil
}

type EcsUserdataDetailsRealRequest struct {
	RegionID   string `json:"regionID,omitempty"`
	InstanceID string `json:"instanceID,omitempty"`
}

type EcsUserdataDetailsRequest struct {
	RegionID   string
	InstanceID string
}

type EcsUserdataDetailsRealResponse struct {
	Userdata string `json:"userdata,omitempty"`
}

type EcsUserdataDetailsResponse struct {
	Userdata string
}
