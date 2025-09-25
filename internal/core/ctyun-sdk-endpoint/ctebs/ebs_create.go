package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// EbsCreateApi 云硬盘开通
// https://www.ctyun.cn/document/10027696/10110700
type EbsCreateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEbsCreateApi(client *ctyunsdk.CtyunClient) *EbsCreateApi {
	return &EbsCreateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ebs/new-ebs",
		},
	}
}

func (this *EbsCreateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EbsCreateRequest) (*EbsCreateResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(ebsCreateRealRequest{
		ClientToken: req.ClientToken,
		DiskName:    req.DiskName,
		DiskMode:    req.DiskMode,
		DiskType:    req.DiskType,
		DiskSize:    req.DiskSize,
		RegionID:    req.RegionId,
		AzName:      req.AzName,
		OnDemand:    req.OnDemand,
		CycleType:   req.CycleType,
		CycleCount:  req.CycleCount,
		ProjectID:   req.ProjectId,
	})
	if err != nil {
		return nil, err
	}

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtebs, builder)
	if err != nil {
		return nil, err
	}

	response := &ebsCreateRealResponse{}
	err = resp.ParseByStandardModelWithCheck(response)
	if err != nil {
		return nil, err
	}

	resources := []EbsCreateResourceResponse{}
	for _, resource := range response.Resources {
		resources = append(resources, EbsCreateResourceResponse{
			OrderId:          resource.OrderID,
			Status:           resource.Status,
			IsMaster:         resource.IsMaster,
			DiskName:         resource.DiskID,
			ResourceType:     resource.ResourceType,
			MasterOrderId:    resource.MasterOrderID,
			UpdateTime:       resource.UpdateTime,
			MasterResourceId: resource.MasterResourceID,
			ItemValue:        resource.ItemValue,
			StartTime:        resource.StartTime,
			CreateTime:       resource.CreateTime,
			DiskId:           resource.DiskID,
		})
	}

	return &EbsCreateResponse{
		MasterResourceStatus: response.MasterResourceStatus,
		RegionId:             response.RegionID,
		MasterOrderId:        response.MasterOrderID,
		MasterResourceId:     response.MasterResourceID,
		MasterOrderNo:        response.MasterOrderNO,
		Resources:            resources,
	}, nil
}

type ebsCreateRealRequest struct {
	ClientToken string `json:"clientToken"`
	DiskName    string `json:"diskName"`
	DiskMode    string `json:"diskMode"`
	DiskType    string `json:"diskType"`
	DiskSize    int64  `json:"diskSize"`
	RegionID    string `json:"regionID"`
	AzName      string `json:"azName"`
	OnDemand    bool   `json:"onDemand"`
	CycleType   string `json:"cycleType"`
	CycleCount  int64  `json:"cycleCount"`
	ProjectID   string `json:"projectID,omitempty"`
}

type ebsCreateRealResponse struct {
	MasterResourceStatus string                         `json:"masterResourceStatus"`
	RegionID             string                         `json:"regionID"`
	MasterOrderID        string                         `json:"masterOrderID"`
	MasterResourceID     string                         `json:"masterResourceID"`
	MasterOrderNO        string                         `json:"masterOrderNO"`
	Resources            []ebsCreateResourceRealRequest `json:"resources"`
}

type ebsCreateResourceRealRequest struct {
	OrderID          string `json:"orderID"`
	Status           int64  `json:"status"`
	IsMaster         bool   `json:"isMaster"`
	DiskName         string `json:"diskName"`
	ResourceType     string `json:"resourceType"`
	MasterOrderID    string `json:"masterOrderID"`
	UpdateTime       int64  `json:"updateTime"`
	MasterResourceID string `json:"masterResourceID"`
	ItemValue        int64  `json:"itemValue"`
	StartTime        int64  `json:"startTime"`
	CreateTime       int64  `json:"createTime"`
	DiskID           string `json:"diskID"`
}

type EbsCreateRequest struct {
	ClientToken string
	DiskName    string
	DiskMode    string
	DiskType    string
	DiskSize    int64
	RegionId    string
	AzName      string
	OnDemand    bool
	CycleType   string
	CycleCount  int64
	ProjectId   string
}

type EbsCreateResponse struct {
	MasterResourceStatus string
	RegionId             string
	MasterOrderId        string
	MasterResourceId     string
	MasterOrderNo        string
	Resources            []EbsCreateResourceResponse
}

type EbsCreateResourceResponse struct {
	OrderId          string
	Status           int64
	IsMaster         bool
	DiskName         string
	ResourceType     string
	MasterOrderId    string
	UpdateTime       int64
	MasterResourceId string
	ItemValue        int64
	StartTime        int64
	CreateTime       int64
	DiskId           string
}
