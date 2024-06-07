package ctvpc

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// SecurityGroupCreateApi 创建安全组
// https://www.ctyun.cn/document/10026755/10040938
type SecurityGroupCreateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewSecurityGroupCreateApi(client *ctyunsdk.CtyunClient) *SecurityGroupCreateApi {
	return &SecurityGroupCreateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/vpc/create-security-group",
		},
	}
}

func (this *SecurityGroupCreateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *SecurityGroupCreateRequest) (*SecurityGroupCreateResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(securityGroupCreateRealRequest{
		ClientToken: req.ClientToken,
		RegionID:    req.RegionId,
		ProjectID:   req.ProjectId,
		VpcID:       req.VpcId,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtvpc, builder)
	if err != nil {
		return nil, err
	}

	result := &securityGroupCreateRealResponse{}
	err = response.ParseByStandardModelWithCheck(result)
	if err != nil {
		return nil, err
	}
	return &SecurityGroupCreateResponse{
		SecurityGroupId: result.SecurityGroupId,
	}, nil
}

type securityGroupCreateRealRequest struct {
	ClientToken string `json:"clientToken"`
	RegionID    string `json:"regionID"`
	ProjectID   string `json:"projectID,omitempty"`
	VpcID       string `json:"vpcID"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type securityGroupCreateRealResponse struct {
	SecurityGroupId string `json:"securityGroupID"`
}

type SecurityGroupCreateRequest struct {
	RegionId    string // 资源池id
	VpcId       string // vpcId
	Name        string // 子网名称 取值范围：1-64，支持数字、字母、中文、_(下划线)、-（中划线）
	Description string // 支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·~！@#￥%……&*（） —— -+={},
	ClientToken string // 客户端存根，用于保证订单幂等性, 长度 1 - 64
	ProjectId   string // 企业项目 ID，默认为0
}

type SecurityGroupCreateResponse struct {
	SecurityGroupId string
}
