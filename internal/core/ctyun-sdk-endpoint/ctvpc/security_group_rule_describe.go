package ctvpc

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// SecurityGroupRuleDescribeApi 查询安全组规则
// https://www.ctyun.cn/document/10026755/10469983
type SecurityGroupRuleDescribeApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewSecurityGroupRuleDescribeApi(client *ctyunsdk.CtyunClient) *SecurityGroupRuleDescribeApi {
	return &SecurityGroupRuleDescribeApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v4/vpc/describe-security-group-rule",
		},
	}
}

func (this *SecurityGroupRuleDescribeApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *SecurityGroupRuleDescribeRequest) (*SecurityGroupRuleDescribeResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	builder.AddParam("regionID", req.RegionId).
		AddParam("securityGroupID", req.SecurityGroupId).
		AddParam("securityGroupRuleID", req.SecurityGroupRuleId)

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtvpc, builder)
	if err != nil {
		return nil, err
	}

	var realResponse securityGroupRuleDescribeRealResponse
	err = response.ParseByStandardModelWithCheck(&realResponse)
	if err != nil {
		return nil, err
	}
	return &SecurityGroupRuleDescribeResponse{
		Direction:       realResponse.Direction,
		Priority:        realResponse.Priority,
		Ethertype:       realResponse.Ethertype,
		Protocol:        realResponse.Protocol,
		Range:           realResponse.Range,
		DestCidrIp:      realResponse.DestCidrIp,
		Description:     realResponse.Description,
		Id:              realResponse.Id,
		SecurityGroupId: realResponse.SecurityGroupID,
		Action:          realResponse.Action,
	}, nil
}

type securityGroupRuleDescribeRealResponse struct {
	Direction       string `json:"direction"`
	Priority        int    `json:"priority"`
	Ethertype       string `json:"ethertype"`
	Protocol        string `json:"protocol"`
	Range           string `json:"range"`
	DestCidrIp      string `json:"destCidrIp"`
	Description     string `json:"description"`
	Id              string `json:"id"`
	SecurityGroupID string `json:"securityGroupID"`
	Action          string `json:"action"`
}

type SecurityGroupRuleDescribeRequest struct {
	RegionId            string // 区域id
	SecurityGroupId     string // 安全组ID
	SecurityGroupRuleId string // 安全组规则id
}

type SecurityGroupRuleDescribeResponse struct {
	Direction       string // 规则方向，出方向则填写egress
	Priority        int    // 优先级，取值范围：[1, 100]，取值越小优先级越大
	Ethertype       string // IP类型，取值范围：IPv4、IPv6
	Protocol        string // 网络协议，取值范围：ANY（任意）、TCP、UDP、ICMP(v4)
	Range           string // 安全组开放的传输层协议相关的源端端口范围
	DestCidrIp      string // 远端地址:0.0.0.0/0
	Description     string // 安全组规则描述信息，满足以下规则： ① 长度0-128字符， ② 支持拉丁字母、中文、数字, 特殊字符 ！@#￥%……&*（）——-+={}《》？：“”【】、；‘'，。、 不能以      http: / https: 开头
	Id              string // id
	SecurityGroupId string // 安全组id
	Action          string // 授权策略，取值范围：accept（允许），drop（拒绝）。
}
