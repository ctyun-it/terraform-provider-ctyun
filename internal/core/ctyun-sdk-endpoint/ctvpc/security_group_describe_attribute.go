package ctvpc

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// SecurityGroupDescribeAttributeApi 查询用户安全组详情
// https://www.ctyun.cn/document/10026755/10040946
type SecurityGroupDescribeAttributeApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewSecurityGroupDescribeAttributeApi(client *ctyunsdk.CtyunClient) *SecurityGroupDescribeAttributeApi {
	return &SecurityGroupDescribeAttributeApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v4/vpc/describe-security-group-attribute",
		},
	}
}

func (this *SecurityGroupDescribeAttributeApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *SecurityGroupDescribeAttributeRequest) (*SecurityGroupDescribeAttributeResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	builder.
		AddParam("regionID", req.RegionId).
		AddParam("securityGroupID", req.SecurityGroupId).
		AddParam("projectID", req.ProjectId).
		AddParam("direction", req.Direction)

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtvpc, builder)
	if err != nil {
		return nil, err
	}

	var realResponse securityGroupDescribeAttributeRealResponse
	err = response.ParseByStandardModelWithCheck(&realResponse)

	var sgrl []SecurityGroupDescribeAttributeSecurityGroupRuleListResponse
	for _, sg := range realResponse.SecurityGroupRuleList {
		sgrl = append(sgrl, SecurityGroupDescribeAttributeSecurityGroupRuleListResponse{
			Direction:       sg.Direction,
			Priority:        sg.Priority,
			Ethertype:       sg.Ethertype,
			Protocol:        sg.Protocol,
			Range:           sg.Range,
			DestCidrIp:      sg.DestCidrIp,
			Description:     sg.Description,
			Origin:          sg.Origin,
			CreateTime:      sg.CreateTime,
			Id:              sg.Id,
			Action:          sg.Action,
			SecurityGroupId: sg.SecurityGroupId,
		})
	}
	return &SecurityGroupDescribeAttributeResponse{
		SecurityGroupName:     realResponse.SecurityGroupName,
		Id:                    realResponse.Id,
		VmNum:                 realResponse.VmNum,
		Origin:                realResponse.Origin,
		VpcName:               realResponse.VpcName,
		VpcId:                 realResponse.VpcId,
		CreationTime:          realResponse.CreationTime,
		Description:           realResponse.Description,
		SecurityGroupRuleList: sgrl,
	}, err
}

type securityGroupDescribeAttributeRealResponse struct {
	SecurityGroupName     string `json:"securityGroupName"`
	Id                    string `json:"id"`
	VmNum                 int    `json:"vmNum"`
	Origin                string `json:"origin"`
	VpcName               string `json:"vpcName"`
	VpcId                 string `json:"vpcID"`
	CreationTime          string `json:"creationTime"`
	Description           string `json:"description"`
	SecurityGroupRuleList []struct {
		Direction       string `json:"direction"`
		Priority        int    `json:"priority"`
		Ethertype       string `json:"ethertype"`
		Protocol        string `json:"protocol"`
		Range           string `json:"range"`
		DestCidrIp      string `json:"destCidrIp"`
		Description     string `json:"description"`
		Origin          string `json:"origin"`
		CreateTime      string `json:"createTime"`
		Id              string `json:"id"`
		Action          string `json:"action"`
		SecurityGroupId string `json:"securityGroupID"`
	} `json:"securityGroupRuleList"`
}

type SecurityGroupDescribeAttributeSecurityGroupRuleListResponse struct {
	Direction       string
	Priority        int
	Ethertype       string
	Protocol        string
	Range           string
	DestCidrIp      string
	Description     string
	Origin          string
	CreateTime      string
	Id              string
	Action          string
	SecurityGroupId string
}

type SecurityGroupDescribeAttributeRequest struct {
	RegionId        string // 区域id
	SecurityGroupId string // 安全组ID
	ProjectId       string // 企业项目 ID，默认为0
	Direction       string // 安全组规则授权方向： egress：安全组出方向 ingress：安全组入方向 all：不区分方向 默认是：all
}

type SecurityGroupDescribeAttributeResponse struct {
	SecurityGroupName     string
	Id                    string
	VmNum                 int
	Origin                string
	VpcName               string
	VpcId                 string
	CreationTime          string
	Description           string
	SecurityGroupRuleList []SecurityGroupDescribeAttributeSecurityGroupRuleListResponse
}
