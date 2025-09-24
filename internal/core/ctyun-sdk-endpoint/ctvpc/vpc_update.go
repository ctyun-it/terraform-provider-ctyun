package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// VpcUpdateApi 修改VPC属性
// https://www.ctyun.cn/document/10026755/10040810
type VpcUpdateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewVpcUpdateApi(client *ctyunsdk.CtyunClient) *VpcUpdateApi {
	return &VpcUpdateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/vpc/update",
		},
	}
}

func (this *VpcUpdateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *VpcUpdateRequest) (*VpcUpdateResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	realRequest := vpcUpdateRealRequest{
		VpcId:       req.VpcId,
		ClientToken: req.ClientToken,
		RegionID:    req.RegionId,
		Name:        req.Name,
		Description: req.Description,
		ProjectID:   req.ProjectId,
	}
	_, err := builder.WriteJson(realRequest)
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
	return &VpcUpdateResponse{}, nil
}

type vpcUpdateRealRequest struct {
	VpcId       string `json:"vpcID"`
	ClientToken string `json:"clientToken"`
	RegionID    string `json:"regionID"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ProjectID   string `json:"projectID,omitempty"`
}
type VpcUpdateRequest struct {
	VpcId       string // 更新的vpcId
	ClientToken string // 客户端存根，用于保证订单幂等性, 长度 1 - 64
	RegionId    string // 资源池id
	Name        string // 虚拟私有云名称
	Description string // 描述
	ProjectId   string // 企业项目 ID，默认为"0"
}

type VpcUpdateResponse struct {
}
