package ctvpc

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// SecurityGroupRuleIngressModifyApi 修改安全组入方向规则
// https://www.ctyun.cn/document/10026755/10040984
type SecurityGroupRuleIngressModifyApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewSecurityGroupRuleIngressModifyApi(client *ctyunsdk.CtyunClient) *SecurityGroupRuleIngressModifyApi {
	return &SecurityGroupRuleIngressModifyApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/vpc/modify-security-group-ingress",
		},
	}
}

func (this *SecurityGroupRuleIngressModifyApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *SecurityGroupRuleIngressModifyRequest) (*SecurityGroupRuleIngressModifyResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&securityGroupRuleIngressModifyRealRequest{
		RegionID:            req.RegionId,
		SecurityGroupID:     req.SecurityGroupId,
		SecurityGroupRuleID: req.SecurityGroupRuleId,
		ClientToken:         req.ClientToken,
		Description:         req.Description,
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
	return &SecurityGroupRuleIngressModifyResponse{}, nil
}

type securityGroupRuleIngressModifyRealRequest struct {
	RegionID            string `json:"regionID"`
	SecurityGroupID     string `json:"securityGroupID"`
	SecurityGroupRuleID string `json:"securityGroupRuleID"`
	ClientToken         string `json:"clientToken"`
	Description         string `json:"description"`
}

type SecurityGroupRuleIngressModifyRequest struct {
	RegionId            string // 资源池ID，请根据查询资源池列表接口返回值进行传参，获取“regionId”参数
	SecurityGroupId     string // 安全组ID
	SecurityGroupRuleId string // 安全规则ID
	ClientToken         string // 客户端存根
	Description         string // 描述
}

type SecurityGroupRuleIngressModifyResponse struct {
}
