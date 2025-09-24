package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// EcsQueryAsyncResultApi 查询一个异步任务的结果
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5543&data=87
type EcsQueryAsyncResultApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEcsQueryAsyncResultApi(client *ctyunsdk.CtyunClient) *EcsQueryAsyncResultApi {
	return &EcsQueryAsyncResultApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v4/ecs/query-async-result",
		},
	}
}

func (this *EcsQueryAsyncResultApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EcsQueryAsyncResultRequest) (*EcsQueryAsyncResultResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&ecsQueryAsyncResultRealRequest{
		RegionID: req.RegionId,
		JobID:    req.JobId,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtecs, builder)
	if err != nil {
		return nil, err
	}

	var realResponse EcsQueryAsyncResultRealResponse
	err = response.ParseByStandardModelWithCheck(&realResponse)
	if err != nil {
		return nil, err
	}
	return &EcsQueryAsyncResultResponse{
		JobStatus: realResponse.JobStatus,
	}, nil
}

type ecsQueryAsyncResultRealRequest struct {
	RegionID string `json:"regionID"`
	JobID    string `json:"jobID"`
}

type EcsQueryAsyncResultRequest struct {
	RegionId string
	JobId    string
}

type EcsQueryAsyncResultRealResponse struct {
	JobStatus int `json:"jobStatus"`
}

type EcsQueryAsyncResultResponse struct {
	JobStatus int
}
