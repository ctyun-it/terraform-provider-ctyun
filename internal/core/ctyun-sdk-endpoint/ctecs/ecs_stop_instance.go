package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// EcsStopInstanceApi 关闭一台云主机
// https://www.ctyun.cn/document/10026730/10106393
type EcsStopInstanceApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEcsStopInstanceApi(client *ctyunsdk.CtyunClient) *EcsStopInstanceApi {
	return &EcsStopInstanceApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ecs/stop-instance",
		},
	}
}

func (this *EcsStopInstanceApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EcsStopInstanceRequest) (*EcsStopInstanceResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&ecsStopInstanceRealRequest{
		RegionID:   req.RegionId,
		InstanceID: req.InstanceId,
		Force:      req.Force,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtecs, builder)
	if err != nil {
		return nil, err
	}

	var realResponse ecsStopInstanceRealResponse
	err = response.ParseByStandardModelWithCheck(&realResponse)
	if err != nil {
		return nil, err
	}
	return &EcsStopInstanceResponse{
		JobId: realResponse.JobID,
	}, nil
}

type ecsStopInstanceRealRequest struct {
	RegionID   string `json:"regionID"`
	InstanceID string `json:"instanceID"`
	Force      bool   `json:"force"`
}

type ecsStopInstanceRealResponse struct {
	JobID string `json:"jobID"`
}

type EcsStopInstanceRequest struct {
	RegionId   string
	InstanceId string
	Force      bool
}

type EcsStopInstanceResponse struct {
	JobId string
}
