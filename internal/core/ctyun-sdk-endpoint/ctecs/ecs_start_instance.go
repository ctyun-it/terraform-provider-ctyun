package ctecs

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// EcsStartInstanceApi 开启一台云主机
// https://www.ctyun.cn/document/10026730/10106397
type EcsStartInstanceApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEcsStartInstanceApi(client *ctyunsdk.CtyunClient) *EcsStartInstanceApi {
	return &EcsStartInstanceApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ecs/start-instance",
		},
	}
}

func (this *EcsStartInstanceApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EcsStartInstanceRequest) (*EcsStartInstanceResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&ecsStartInstanceRealRequest{
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

	var realResponse EcsStartInstanceRealResponse
	err = response.ParseByStandardModelWithCheck(&realResponse)
	if err != nil {
		return nil, err
	}
	return &EcsStartInstanceResponse{
		JobId: realResponse.JobID,
	}, nil
}

type ecsStartInstanceRealRequest struct {
	RegionID   string `json:"regionID"`
	InstanceID string `json:"instanceID"`
	Force      bool   `json:"force"`
}

type EcsStartInstanceRequest struct {
	RegionId   string
	InstanceId string
	Force      bool
}

type EcsStartInstanceRealResponse struct {
	JobID string `json:"jobID"`
}

type EcsStartInstanceResponse struct {
	JobId string
}
