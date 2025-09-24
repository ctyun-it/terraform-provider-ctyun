package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// SecurityGroupModifyAttributionApi 修改安全组
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=6319&data=94
type SecurityGroupModifyAttributionApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewSecurityGroupModifyAttributionApi(client *ctyunsdk.CtyunClient) *SecurityGroupModifyAttributionApi {
	return &SecurityGroupModifyAttributionApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/vpc/modify-security-group-attribute",
		},
	}
}

func (this *SecurityGroupModifyAttributionApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *SecurityGroupModifyAttributionRequest) (*SecurityGroupModifyAttributionResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(securityGroupModifyAttributionRealRequest{
		ClientToken:     req.ClientToken,
		RegionID:        req.RegionId,
		ProjectID:       req.ProjectId,
		Name:            req.Name,
		Description:     req.Description,
		Enabled:         req.Enabled,
		SecurityGroupID: req.SecurityGroupId,
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
	return &SecurityGroupModifyAttributionResponse{}, nil
}

type securityGroupModifyAttributionRealRequest struct {
	ClientToken     string `json:"clientToken"`
	RegionID        string `json:"regionID"`
	ProjectID       string `json:"projectID,omitempty"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Enabled         bool   `json:"enabled"`
	SecurityGroupID string `json:"securityGroupID"`
}

type SecurityGroupModifyAttributionRequest struct {
	ClientToken     string // 客户端存根，用于保证订单幂等性, 长度 1 - 64
	RegionId        string // 资源池id
	ProjectId       string // 企业项目 ID，默认为0
	Name            string // 子网名称 取值范围：1-64，支持数字、字母、中文、_(下划线)、-（中划线）
	Description     string // 支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·~！@#￥%……&*（） —— -+={},
	Enabled         bool   // 开启安全组 / 关闭安全组
	SecurityGroupId string // 安全组ID
}

type SecurityGroupModifyAttributionResponse struct {
}
