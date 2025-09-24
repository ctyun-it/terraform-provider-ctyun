package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// EcsJoinSecurityGroupApi 绑定安全组
// https://www.ctyun.cn/document/10026730/10040193
type EcsJoinSecurityGroupApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEcsJoinSecurityGroupApi(client *ctyunsdk.CtyunClient) *EcsJoinSecurityGroupApi {
	return &EcsJoinSecurityGroupApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ecs/vpc/join-security-group",
		},
	}
}

func (this *EcsJoinSecurityGroupApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EcsJoinSecurityGroupRequest) (*EcsJoinSecurityGroupResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&ecsJoinSecurityGroupRealRequest{
		RegionID:           req.RegionId,
		SecurityGroupID:    req.SecurityGroupId,
		InstanceID:         req.InstanceId,
		NetworkInterfaceID: req.NetworkInterfaceId,
		Action:             req.Action,
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
	return &EcsJoinSecurityGroupResponse{}, nil
}

type ecsJoinSecurityGroupRealRequest struct {
	RegionID           string `json:"regionID,omitempty"`
	SecurityGroupID    string `json:"securityGroupID,omitempty"`
	InstanceID         string `json:"instanceID,omitempty"`
	NetworkInterfaceID string `json:"networkInterfaceID,omitempty"`
	Action             string `json:"action,omitempty"`
}

type EcsJoinSecurityGroupRequest struct {
	RegionId           string
	SecurityGroupId    string
	InstanceId         string
	NetworkInterfaceId string
	Action             string
}

type EcsJoinSecurityGroupResponse struct {
}
