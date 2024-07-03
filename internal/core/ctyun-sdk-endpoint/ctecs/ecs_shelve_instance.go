package ctecs

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// EcsShelveInstanceApi 节省关机一台云主机
// https://www.ctyun.cn/document/10026730/10597755
type EcsShelveInstanceApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEcsShelveInstanceApi(client *ctyunsdk.CtyunClient) *EcsShelveInstanceApi {
	return &EcsShelveInstanceApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ecs/shelve-instance",
		},
	}
}

func (this *EcsShelveInstanceApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EcsShelveInstanceRequest) (*EcsShelveInstanceResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)

	_, err := builder.WriteJson(&EcsShelveInstanceRealRequest{
		RegionID:   req.RegionID,
		InstanceID: req.InstanceID,
	})

	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtecs, builder)
	if err != nil {
		return nil, err
	}

	var realResponse EcsShelveInstanceRealResponse
	err = response.ParseByStandardModelWithCheck(&realResponse)
	if err != nil {
		return nil, err
	}

	return &EcsShelveInstanceResponse{
		JobID: realResponse.JobID,
	}, nil
}

type EcsShelveInstanceRealRequest struct {
	RegionID   string `json:"regionID,omitempty"`
	InstanceID string `json:"instanceID,omitempty"`
}

type EcsShelveInstanceRequest struct {
	RegionID   string
	InstanceID string
}

type EcsShelveInstanceRealResponse struct {
	JobID string `json:"jobID,omitempty"`
}

type EcsShelveInstanceResponse struct {
	JobID string
}
