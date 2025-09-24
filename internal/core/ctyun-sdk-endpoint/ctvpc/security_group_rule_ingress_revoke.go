package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// SecurityGroupRuleIngressRevokeApi 删除安全组入方向规则
// https://www.ctyun.cn/document/10026755/10040977
type SecurityGroupRuleIngressRevokeApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewSecurityGroupRuleIngressRevokeApi(client *ctyunsdk.CtyunClient) *SecurityGroupRuleIngressRevokeApi {
	return &SecurityGroupRuleIngressRevokeApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/vpc/revoke-security-group-ingress",
		},
	}
}

func (this *SecurityGroupRuleIngressRevokeApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *SecurityGroupRuleIngressRevokeRequest) (*SecurityGroupRuleIngressRevokeResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&securityGroupRuleIngressRevokeRealRequest{
		RegionID:            req.RegionId,
		SecurityGroupID:     req.SecurityGroupId,
		SecurityGroupRuleID: req.SecurityGroupRuleId,
		ClientToken:         req.ClientToken,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtvpc, builder)
	if err != nil {
		return nil, err
	}

	err = response.ParseByStandardModelWithCheck(nil)
	if err != nil {
		return nil, err
	}
	return &SecurityGroupRuleIngressRevokeResponse{}, nil
}

type securityGroupRuleIngressRevokeRealRequest struct {
	RegionID            string `json:"regionID"`
	SecurityGroupID     string `json:"securityGroupID"`
	SecurityGroupRuleID string `json:"securityGroupRuleID"`
	ClientToken         string `json:"clientToken"`
}

type SecurityGroupRuleIngressRevokeRequest struct {
	RegionId            string // 资源池ID，请根据查询资源池列表接口返回值进行传参，获取“regionId”参数
	SecurityGroupId     string // 安全组ID
	SecurityGroupRuleId string // 安全规则ID
	ClientToken         string // 客户端存根
}

type SecurityGroupRuleIngressRevokeResponse struct {
}
