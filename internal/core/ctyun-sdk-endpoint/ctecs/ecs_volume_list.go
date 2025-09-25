package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// EcsVolumeListApi 查询云主机的云硬盘列表
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8290&data=87
type EcsVolumeListApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEcsVolumeListApi(client *ctyunsdk.CtyunClient) *EcsVolumeListApi {
	return &EcsVolumeListApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ecs/volume/list",
		},
	}
}

func (this *EcsVolumeListApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EcsVolumeListRequest) (*EcsVolumeListResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&ecsVolumeListRealRequest{
		RegionID:   req.RegionId,
		InstanceID: req.InstanceId,
		PageNo:     req.PageNo,
		PageSize:   req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtecs, builder)
	if err != nil {
		return nil, err
	}

	var realResponse ecsVolumeListRealResponse
	err = response.ParseByStandardModelWithCheck(&realResponse)
	if err != nil {
		return nil, err
	}
	var results []EcsVolumeListResultsResponse
	for _, result := range realResponse.Results {
		results = append(results, EcsVolumeListResultsResponse{
			DiskType:     result.DiskType,
			IsEncrypt:    result.IsEncrypt,
			DiskSize:     result.DiskSize,
			DiskMode:     result.DiskMode,
			DiskId:       result.DiskID,
			DiskDataType: result.DiskDataType,
		})
	}
	return &EcsVolumeListResponse{
		CurrentCount: realResponse.CurrentCount,
		TotalCount:   realResponse.TotalCount,
		TotalPage:    realResponse.TotalPage,
		Results:      results,
	}, nil
}

type ecsVolumeListRealRequest struct {
	RegionID   string `json:"regionID,omitempty"`
	InstanceID string `json:"instanceID,omitempty"`
	PageNo     int    `json:"securityGroupID,omitempty"`
	PageSize   int    `json:"networkInterfaceID,omitempty"`
}

type EcsVolumeListRequest struct {
	RegionId   string
	InstanceId string
	PageNo     int
	PageSize   int
}

type ecsVolumeListRealResponse struct {
	CurrentCount int `json:"currentCount"`
	TotalCount   int `json:"totalCount"`
	TotalPage    int `json:"totalPage"`
	Results      []struct {
		DiskType     string `json:"diskType"`
		IsEncrypt    bool   `json:"isEncrypt"`
		DiskSize     int    `json:"diskSize"`
		DiskMode     string `json:"diskMode"`
		DiskID       string `json:"diskID"`
		DiskDataType string `json:"diskDataType"`
	} `json:"results"`
}

type EcsVolumeListResultsResponse struct {
	DiskType     string
	IsEncrypt    bool
	DiskSize     int
	DiskMode     string
	DiskId       string
	DiskDataType string
}

type EcsVolumeListResponse struct {
	CurrentCount int
	TotalCount   int
	TotalPage    int
	Results      []EcsVolumeListResultsResponse
}
