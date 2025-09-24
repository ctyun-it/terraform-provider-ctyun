package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// SecurityGroupRuleEgressCreateApi 创建安全组出向规则
// https://www.ctyun.cn/document/10026730/10040197
type SecurityGroupRuleEgressCreateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewSecurityGroupRuleEgressCreateApi(client *ctyunsdk.CtyunClient) *SecurityGroupRuleEgressCreateApi {
	return &SecurityGroupRuleEgressCreateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ecs/vpc/create-security-group-egress",
		},
	}
}

func (this *SecurityGroupRuleEgressCreateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *SecurityGroupRuleEgressCreateRequest) (*SecurityGroupRuleEgressCreateResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	var securityGroupRules []securityGroupRuleEgressCreateSecurityGroupRulesRealRequest
	for _, rule := range req.SecurityGroupRules {
		securityGroupRules = append(securityGroupRules, securityGroupRuleEgressCreateSecurityGroupRulesRealRequest{
			Direction:   rule.Direction,
			Action:      rule.Action,
			Priority:    rule.Priority,
			Protocol:    rule.Protocol,
			Ethertype:   rule.Ethertype,
			DestCidrIp:  rule.DestCidrIp,
			Description: rule.Description,
			Range:       rule.Range,
		})
	}
	_, err := builder.WriteJson(&securityGroupRuleEgressCreateRealRequest{
		RegionID:           req.RegionId,
		SecurityGroupID:    req.SecurityGroupId,
		ClientToken:        req.ClientToken,
		SecurityGroupRules: securityGroupRules,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtecs, builder)
	if err != nil {
		return nil, err
	}

	var realResponse securityGroupRuleEgressCreateRealResponse
	err = response.ParseByStandardModelWithCheck(&realResponse)
	if err != nil {
		return nil, err
	}
	return &SecurityGroupRuleEgressCreateResponse{
		SgRuleIds: realResponse.SgRuleIDs,
	}, nil
}

type securityGroupRuleEgressCreateSecurityGroupRulesRealRequest struct {
	Direction   string `json:"direction"`
	Action      string `json:"action"`
	Priority    int    `json:"priority"`
	Protocol    string `json:"protocol"`
	Ethertype   string `json:"ethertype"`
	DestCidrIp  string `json:"destCidrIp"`
	Description string `json:"description"`
	Range       string `json:"range"`
}

type securityGroupRuleEgressCreateRealRequest struct {
	RegionID           string                                                       `json:"regionID"`
	SecurityGroupID    string                                                       `json:"securityGroupID"`
	ClientToken        string                                                       `json:"clientToken"`
	SecurityGroupRules []securityGroupRuleEgressCreateSecurityGroupRulesRealRequest `json:"securityGroupRules"`
}

type securityGroupRuleEgressCreateRealResponse struct {
	SgRuleIDs []string `json:"sgRuleIDs"`
}

type SecurityGroupRuleEgressCreateSecurityGroupRulesRequest struct {
	Direction   string // 规则方向，出方向则填写egress
	Action      string // 授权策略，取值范围：accept（允许），drop（拒绝）。
	Priority    int    // 优先级，取值范围：[1, 100]，取值越小优先级越大
	Protocol    string // 网络协议，取值范围：ANY（任意）、TCP、UDP、ICMP(v4)
	Ethertype   string // IP类型，取值范围：IPv4、IPv6
	DestCidrIp  string // 远端地址:0.0.0.0/0
	Description string // 安全组规则描述信息，满足以下规则： ① 长度0-128字符， ② 支持拉丁字母、中文、数字, 特殊字符 ！@#￥%……&*（）——-+={}《》？：“”【】、；‘'，。、 不能以      http: / https: 开头
	Range       string // 安全组开放的传输层协议相关的源端端口范围
}

type SecurityGroupRuleEgressCreateRequest struct {
	RegionId           string                                                   // 区域id
	SecurityGroupId    string                                                   // 安全组ID。
	ClientToken        string                                                   // 客户端存根
	SecurityGroupRules []SecurityGroupRuleEgressCreateSecurityGroupRulesRequest // 规则信息
}

type SecurityGroupRuleEgressCreateResponse struct {
	SgRuleIds []string
}
