package ctvpc

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// SecurityGroupDeleteApi 删除安全组
// https://www.ctyun.cn/document/10026755/10040966
type SecurityGroupDeleteApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewSecurityGroupDeleteApi(client *ctyunsdk.CtyunClient) *SecurityGroupDeleteApi {
	return &SecurityGroupDeleteApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/vpc/delete-security-group",
		},
	}
}

func (this *SecurityGroupDeleteApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *SecurityGroupDeleteRequest) (*SecurityGroupDeleteResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(securityGroupDeleteRealRequest{
		ClientToken:     req.ClientToken,
		RegionID:        req.RegionId,
		ProjectID:       req.ProjectId,
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
	return &SecurityGroupDeleteResponse{}, err
}

type securityGroupDeleteRealRequest struct {
	ClientToken     string `json:"clientToken"`
	RegionID        string `json:"regionID"`
	ProjectID       string `json:"projectID,omitempty"`
	SecurityGroupID string `json:"securityGroupID"`
}

type SecurityGroupDeleteRequest struct {
	ClientToken     string // 客户端存根，用于保证订单幂等性, 长度 1 - 64
	RegionId        string // 资源池id
	ProjectId       string // 企业项目 ID，默认为0
	SecurityGroupId string // 安全组ID
}

type SecurityGroupDeleteResponse struct {
}
