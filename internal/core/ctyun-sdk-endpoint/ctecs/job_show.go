package ctecs

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// JobShowApi 通用任务状态查询
// https://www.ctyun.cn/document/10026730/10040594
type JobShowApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewJobShowApi(client *ctyunsdk.CtyunClient) *JobShowApi {
	return &JobShowApi{
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v4/job/info",
		},
		client: client,
	}
}

func (this *JobShowApi) Do(ctx context.Context, credential ctyunsdk.Credential, t *JobShowRequest) (*JobShowResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	builder.AddParam("jobID", t.JobId)
	builder.AddParam("regionID", t.RegionId)
	resp, requestError := this.client.RequestToEndpoint(ctx, EndpointNameCtecs, builder)
	if requestError != nil {
		return nil, requestError
	}
	response := &jobShowRealResponse{}
	err := resp.ParseByStandardModelWithCheck(response)
	if err != nil {
		return nil, err
	}
	return &JobShowResponse{
		Status:     response.Status,
		ResourceId: response.ResourceId,
		Fields:     response.Fields,
		JobStatus:  response.JobStatus,
		JobId:      response.JobID,
	}, nil
}

type JobShowRequest struct {
	RegionId string
	JobId    string
}

type jobShowRealResponse struct {
	Status     int               `json:"status"`
	ResourceId string            `json:"resourceId"`
	Fields     map[string]string `json:"fields"`
	JobStatus  string            `json:"jobStatus"`
	JobID      string            `json:"jobID"`
}

type JobShowResponse struct {
	Status     int
	ResourceId string
	Fields     map[string]string
	JobStatus  string
	JobId      string
}
