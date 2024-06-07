package ctvpc

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// SecurityGroupRuleEgressRevokeApi 删除安全组出方向规则
// https://www.ctyun.cn/document/10026755/10040970
type SecurityGroupRuleEgressRevokeApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewSecurityGroupRuleEgressRevokeApi(client *ctyunsdk.CtyunClient) *SecurityGroupRuleEgressRevokeApi {
	return &SecurityGroupRuleEgressRevokeApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/vpc/revoke-security-group-egress",
		},
	}
}

func (this *SecurityGroupRuleEgressRevokeApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *SecurityGroupRuleEgressRevokeRequest) (*SecurityGroupRuleEgressRevokeResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&securityGroupRuleEgressRevokeRealRequest{
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
	return &SecurityGroupRuleEgressRevokeResponse{}, nil
}

type securityGroupRuleEgressRevokeRealRequest struct {
	RegionID            string `json:"regionID"`
	SecurityGroupID     string `json:"securityGroupID"`
	SecurityGroupRuleID string `json:"securityGroupRuleID"`
	ClientToken         string `json:"clientToken"`
}

type SecurityGroupRuleEgressRevokeRequest struct {
	RegionId            string // 资源池ID，请根据查询资源池列表接口返回值进行传参，获取“regionId”参数
	SecurityGroupId     string // 安全组ID
	SecurityGroupRuleId string // 安全规则ID
	ClientToken         string // 客户端存根
}

type SecurityGroupRuleEgressRevokeResponse struct {
}
