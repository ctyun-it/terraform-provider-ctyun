package ctecs

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// EcsInstanceStatusListApi 获取多台云主机状态
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8308&data=87
type EcsInstanceStatusListApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEcsInstanceStatusListApi(client *ctyunsdk.CtyunClient) *EcsInstanceStatusListApi {
	return &EcsInstanceStatusListApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ecs/instance-status-list",
		},
	}
}

func (this *EcsInstanceStatusListApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EcsInstanceStatusListRequest) (*EcsInstanceStatusListResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&ecsInstanceStatusListRealRequest{
		RegionId:       req.RegionId,
		AzName:         req.AzName,
		InstanceIDList: req.InstanceIdList,
		PageNo:         req.PageNo,
		PageSize:       req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtecs, builder)
	if err != nil {
		return nil, err
	}

	var realResponse ecsInstanceStatusListRealResponse
	err = response.ParseByStandardModelWithCheck(&realResponse)
	if err != nil {
		return nil, err
	}
	var ecsInstanceStatusListStatusListResponse []EcsInstanceStatusListStatusListResponse
	for _, listRealResponse := range realResponse.StatusList {
		ecsInstanceStatusListStatusListResponse = append(ecsInstanceStatusListStatusListResponse, EcsInstanceStatusListStatusListResponse{
			InstanceId:     listRealResponse.InstanceID,
			InstanceStatus: listRealResponse.InstanceStatus,
		})
	}
	return &EcsInstanceStatusListResponse{
		CurrentCount: realResponse.CurrentCount,
		TotalCount:   realResponse.TotalCount,
		TotalPage:    realResponse.TotalPage,
		StatusList:   ecsInstanceStatusListStatusListResponse,
	}, nil
}

type ecsInstanceStatusListRealRequest struct {
	RegionId       string `json:"regionID"`
	AzName         string `json:"azName"`
	InstanceIDList string `json:"instanceIDList"`
	PageNo         int    `json:"pageNo"`
	PageSize       int    `json:"pageSize"`
}

type ecsInstanceStatusListRealResponse struct {
	CurrentCount int `json:"currentCount"`
	TotalCount   int `json:"totalCount"`
	TotalPage    int `json:"totalPage"`
	StatusList   []struct {
		InstanceID     string `json:"instanceID"`
		InstanceStatus string `json:"instanceStatus"`
	} `json:"statusList"`
}

type EcsInstanceStatusListRequest struct {
	RegionId       string
	AzName         string
	InstanceIdList string
	PageNo         int
	PageSize       int
}

type EcsInstanceStatusListResponse struct {
	CurrentCount int
	TotalCount   int
	TotalPage    int
	StatusList   []EcsInstanceStatusListStatusListResponse
}

type EcsInstanceStatusListStatusListResponse struct {
	InstanceId     string
	InstanceStatus string
}
