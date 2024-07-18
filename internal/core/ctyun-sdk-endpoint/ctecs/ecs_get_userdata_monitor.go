package ctecs

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// EcsGetUserdataMonitorApi
type EcsGetUserdataMonitorApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEcsGetUserdataMonitorApi(client *ctyunsdk.CtyunClient) *EcsGetUserdataMonitorApi {
	return &EcsGetUserdataMonitorApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v4/ecs/get-userdata-monitor",
		},
	}
}

func (this *EcsGetUserdataMonitorApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EcsGetUserdataMonitorRequest) (*EcsGetUserdataMonitorResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	builder.AddParam("regionID", req.RegionID)
	builder.AddParam("instanceIDs", req.InstanceIDs)

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtecs, builder)
	if err != nil {
		return nil, err
	}

	var realResponse EcsGetUserdataMonitorRealResponse
	err = response.ParseByStandardModelWithCheck(&realResponse)
	if err != nil {
		return nil, err
	}

	var results []EcsGetUserdataMonitorResultsResponse
	for _, res := range realResponse.Results {
		results = append(results, EcsGetUserdataMonitorResultsResponse{
			InstanceID:     res.InstanceID,
			UserData:       res.UserData,
			MonitorService: res.MonitorService,
		})
	}

	return &EcsGetUserdataMonitorResponse{
		Results: results,
	}, nil
}

type EcsGetUserdataMonitorRealRequest struct {
	RegionID    string `json:"regionID,omitempty"`
	InstanceIDs string `json:"instanceIDs,omitempty"`
}

type EcsGetUserdataMonitorRequest struct {
	RegionID    string
	InstanceIDs string
}

type EcsGetUserdataMonitorResultsRealResponse struct {
	InstanceID     string `json:"instanceID,omitempty"`
	UserData       string `json:"userData,omitempty"`
	MonitorService bool   `json:"monitorService,omitempty"`
}

type EcsGetUserdataMonitorRealResponse struct {
	Results []EcsGetUserdataMonitorResultsRealResponse `json:"results,omitempty"`
}

type EcsGetUserdataMonitorResultsResponse struct {
	InstanceID     string
	UserData       string
	MonitorService bool
}

type EcsGetUserdataMonitorResponse struct {
	Results []EcsGetUserdataMonitorResultsResponse
}
