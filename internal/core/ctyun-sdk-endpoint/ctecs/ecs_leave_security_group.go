package ctecs

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// EcsLeaveSecurityGroupApi  绑定安全组
// https://www.ctyun.cn/document/10026730/10040193
type EcsLeaveSecurityGroupApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEcsLeaveSecurityGroupApi(client *ctyunsdk.CtyunClient) *EcsLeaveSecurityGroupApi {
	return &EcsLeaveSecurityGroupApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ecs/vpc/leave-security-group",
		},
	}
}

func (this *EcsLeaveSecurityGroupApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EcsLeaveSecurityGroupRequest) (*EcsLeaveSecurityGroupResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&ecsLeaveSecurityGroupRealRequest{
		RegionID:        req.RegionId,
		SecurityGroupID: req.SecurityGroupId,
		InstanceID:      req.InstanceId,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtecs, builder)
	if err != nil {
		return nil, err
	}

	err = response.ParseByStandardModelWithCheck(nil)
	if err != nil {
		return nil, err
	}
	return &EcsLeaveSecurityGroupResponse{}, nil
}

type ecsLeaveSecurityGroupRealRequest struct {
	RegionID        string `json:"regionID,omitempty"`
	SecurityGroupID string `json:"securityGroupID,omitempty"`
	InstanceID      string `json:"instanceID,omitempty"`
}

type EcsLeaveSecurityGroupRequest struct {
	RegionId        string
	SecurityGroupId string
	InstanceId      string
}

type EcsLeaveSecurityGroupResponse struct {
}
