package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbListLoadBalancerApi
/* 查看负载均衡实例列表
 */type CtelbListLoadBalancerApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbListLoadBalancerApi(client *core.CtyunClient) *CtelbListLoadBalancerApi {
	return &CtelbListLoadBalancerApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/elb/list-loadbalancer",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbListLoadBalancerApi) Do(ctx context.Context, credential core.Credential, req *CtelbListLoadBalancerRequest) (*CtelbListLoadBalancerResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.IDs != "" {
		ctReq.AddParam("IDs", req.IDs)
	}
	if req.ResourceType != "" {
		ctReq.AddParam("resourceType", req.ResourceType)
	}
	if req.Name != "" {
		ctReq.AddParam("name", req.Name)
	}
	if req.SubnetID != "" {
		ctReq.AddParam("subnetID", req.SubnetID)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtelbListLoadBalancerResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbListLoadBalancerRequest struct {
	RegionID     string /*  区域ID  */
	IDs          string /*  负载均衡ID列表，以,分隔  */
	ResourceType string /*  资源类型。internal：内网负载均衡，external：公网负载均衡  */
	Name         string /*  名称  */
	SubnetID     string /*  子网ID  */
}

type CtelbListLoadBalancerResponse struct {
	StatusCode  int32                                     `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                    `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                    `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                    `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtelbListLoadBalancerReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       string                                    `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbListLoadBalancerReturnObjResponse struct {
	RegionID         string                                        `json:"regionID,omitempty"`         /*  区域ID  */
	AzName           string                                        `json:"azName,omitempty"`           /*  可用区名称  */
	ID               string                                        `json:"ID,omitempty"`               /*  负载均衡ID  */
	ProjectID        string                                        `json:"projectID,omitempty"`        /*  项目ID  */
	Name             string                                        `json:"name,omitempty"`             /*  名称  */
	Description      string                                        `json:"description,omitempty"`      /*  描述  */
	VpcID            string                                        `json:"vpcID,omitempty"`            /*  VPC ID  */
	SubnetID         string                                        `json:"subnetID,omitempty"`         /*  子网ID  */
	PortID           string                                        `json:"portID,omitempty"`           /*  负载均衡实例默认创建port ID  */
	PrivateIpAddress string                                        `json:"privateIpAddress,omitempty"` /*  负载均衡实例的内网VIP  */
	Ipv6Address      string                                        `json:"ipv6Address,omitempty"`      /*  负载均衡实例的IPv6地址  */
	EipInfo          []*CtelbListLoadBalancerReturnObjEipInfoModel `json:"eipInfo"`                    /*  弹性公网IP信息  */
	SlaName          string                                        `json:"slaName,omitempty"`          /*  规格名称  */
	DeleteProtection *bool                                         `json:"deleteProtection"`           /*  删除保护。开启，不开启  */
	AdminStatus      string                                        `json:"adminStatus,omitempty"`      /*  管理状态: DOWN / ACTIVE  */
	Status           string                                        `json:"status,omitempty"`           /*  负载均衡状态: DOWN / ACTIVE  */
	ResourceType     string                                        `json:"resourceType,omitempty"`     /*  负载均衡类型: external / internal  */
	CreatedTime      string                                        `json:"createdTime,omitempty"`      /*  创建时间，为UTC格式  */
	UpdatedTime      string                                        `json:"updatedTime,omitempty"`      /*  更新时间，为UTC格式  */
}

type CtelbListLoadBalancerReturnObjEipInfoModel struct {
	ResourceID  string `json:"resourceID,omitempty"` /*  计费类资源ID  */
	EipID       string `json:"eipID,omitempty"`      /*  弹性公网IP的ID  */
	Bandwidth   int32  `json:"bandwidth,omitempty"`  /*  弹性公网IP的带宽  */
	IsTalkOrder *bool  `json:"isTalkOrder"`          /*  是否按需资源  */
}
